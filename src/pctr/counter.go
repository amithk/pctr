package pctr

import "encoding/binary"
import "io"
import "os"
import "path/filepath"
import "sync"

// Persistent counter optimized for single process, multithreaded scenarios.
// In other words, all accesses to the object of PersistentCounter are thread safe.
// But if the same counter file is accessed by creating multiple objects of this
// structure, the results can be inconsistent.
type PersistentCounter struct {
	counterid string
	spath     string
	Value     uint64
	deleted   bool
	f         *os.File
	mutx      *sync.Mutex
}

func NewPersistentCounter(spath, counterid string) (*PersistentCounter, error) {
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

	pc.Value = val
	return pc, nil
}

// IncrementValue returns the updated value of the counter. This make GetValue
// API unnecessary. For all practical purposes, NewPersistentCounter,
// IncrementValue and DeleteCounter APIs should be enough.
func (pc *PersistentCounter) IncrementValue(incr uint64) (uint64, error) {
	pc.mutx.Lock()
	defer pc.mutx.Unlock()
	newVal := pc.Value + incr
	buf := serializeUint64(newVal)
	err := WriteFile(pc.f, buf)
	if err != nil {
		return 0, err
	}

	pc.Value = newVal
	return newVal, err
}

func (pc *PersistentCounter) DeleteCounter() error {
	return nil
}

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

func (pc *PersistentCounter) IsDeleted() bool {
	return false
}

// Persistent counter optimized for mutliprocess scenarios.
type PersistentCounterMultiproc struct {
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
