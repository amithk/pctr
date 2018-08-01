package pctr

import "encoding/binary"
import "errors"
import "io"
import "os"
import "path/filepath"
import "sync"
import "sync/atomic"

var counterMap map[string]*PersistentCounter = make(map[string]*PersistentCounter)
var counterMapLock *sync.Mutex = &sync.Mutex{}

var ErrCounterNotInitialized error = errors.New("Counter is not initialized")
var ErrCounterDeleted error = errors.New("Access deleted counter error")

// Persistent counter optimized for single process, multithreaded scenarios.
// In other words, all accesses to the object of PersistentCounter are thread safe.
// But if the same counter file is accessed by creating multiple objects of this
// structure, the results can be inconsistent.
type PersistentCounter struct {
	counterid   string // has to be unique
	spath       string
	curMax      uint64
	curVal      uint64
	deleted     bool
	f           *os.File
	mutx        *sync.Mutex
	initialized bool
}

func NewPersistentCounter(spath, counterid string) (*PersistentCounter, error) {
	if _, ok := counterMap[counterid]; ok {
		return nil, errors.New("Attempt to create second object for same counter")
	}

	counterMapLock.Lock()
	defer counterMapLock.Unlock()

	if _, ok := counterMap[counterid]; ok {
		return nil, errors.New("Attempt to create second object for same counter")
	}

	pc := &PersistentCounter{
		counterid: counterid,
		spath:     spath,
		mutx:      &sync.Mutex{},
	}
	f, err := OpenFile(filepath.Join(spath, counterid))
	if err != nil {
		return nil, err
	}
	pc.f = f

	val, err1 := pc.GetValue()
	if err1 != nil {
		return nil, err1
	}

	atomic.StoreUint64(&pc.curMax, val)
	atomic.StoreUint64(&pc.curVal, val)
	counterMap[counterid] = pc
	pc.initialized = true
	return pc, nil
}

// IncrementValue returns the updated value of the counter. This make GetValue
// API unnecessary. For all practical purposes, NewPersistentCounter,
// IncrementValue and DeleteCounter APIs should be enough.
func (pc *PersistentCounter) IncrementValue(incr uint64) (uint64, error) {
	if pc.deleted {
		// This error is primarily used as a circuit breaker.
		return 0, ErrCounterDeleted
	}

	if !pc.initialized {
		return 0, ErrCounterNotInitialized
	}

	pc.mutx.Lock()
	defer pc.mutx.Unlock()
	lm := atomic.LoadUint64(&pc.curMax)
	newVal := lm + incr
	buf := serializeUint64(newVal)
	err := WriteFile(pc.f, buf)
	if err != nil {
		return 0, err
	}

	atomic.StoreUint64(&pc.curMax, newVal)
	return newVal, err
}

func (pc *PersistentCounter) DeleteCounter() error {
	pc.deleted = true
	pc.mutx.Lock()
	defer pc.mutx.Unlock()
	return DeleteFile(filepath.Join(pc.spath, pc.counterid))
}

func (pc *PersistentCounter) GetNext() (uint64, error) {
	if pc.deleted {
		// This error is primarily used as a circuit breaker.
		return 0, ErrCounterDeleted
	}

	if !pc.initialized {
		return 0, ErrCounterNotInitialized
	}

	for {
		lv := atomic.LoadUint64(&pc.curVal)
		lm := atomic.LoadUint64(&pc.curMax)
		if lv < lm {
			next := atomic.AddUint64(&pc.curVal, 1)
			if next > lm {
				continue
			}
			return next, nil
		} else {
			// Keep batchsize hard-coded, for now.
			_, err := pc.IncrementValue(64)
			if err != nil {
				return 0, err
			}
			continue
		}
	}
}

// Note that GetValue is expensive operation. It does uncached read from file
// everytime. GetNext should be called to allocate new counter values.
func (pc *PersistentCounter) GetValue() (uint64, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	err := ReadFile(pc.f, buf)
	if err == io.EOF {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return deserializeUint64(buf), nil
}

func serializeUint64(val uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, val)
	return buf
}

func deserializeUint64(buf []byte) uint64 {
	val, _ := binary.Uvarint(buf)
	return val
}
