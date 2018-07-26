package pctr

import "os"
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

func NewPersistentCounter(counterid, spath string) (*PersistentCounter, error) {
	pc := &PersistentCounter{
		counterid: counterid,
		spath:     spath,
		mutx:      &sync.Mutex{},
	}
	f, err := CreateNew(spath)
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
	newVal := pc.Value + incr
	err := WriteFile(pc.f, newVal)
	pc.Value = newVal
	pc.mutx.Unlock()
	return newVal, err
}

func (pc *PersistentCounter) DeleteCounter() error {
	return nil
}

func (pc *PersistentCounter) GetValue() (uint64, error) {
	return ReadFile(pc.f)
}

func (pc *PersistentCounter) IsDeleted() bool {
	return false
}

// Persistent counter optimized for mutliprocess scenarios.
type PersistentCounterMultiproc struct {
}
