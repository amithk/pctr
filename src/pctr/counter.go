package pctr

type CounterValue struct {
	val uint64
}

type PersistentCounter struct {
	counterid string
	spath     string
	Value     CounterValue
	deleted   bool
	fileops   *FileOps
}

func NewPersistentCounter(counterid, spath string) (*PersistentCounter, error) {
	pc := &PersistentCounter{
		counterid: counterid,
		spath:     spath,
	}
	pc.fileops = &FileOps{}
	pc.fileops.CreateNew(spath)
}

func (pc *PersistentCounter) IncrementValue(incr uint64) (uint64, error) {
}

func (pc *PersistentCounter) DeleteCounter() error {
}

func (pc *PersistentCounter) GetValue() (uint64, error) {
}

func (pc *PersistentCounter) IsDeleted() bool {
}
