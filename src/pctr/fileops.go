package pctr

type FileOps struct {
}

func (fo *FileOps) CreateNew(fpath string) error {
}

func (fo *FileOps) Move(oldPath, newPath string) error {
}

func (fo *FileOps) Read(fpath string) (uint64, error) {
}

func (fo *FileOps) Write(fpath string, val uint64) error {
}

func (fo *FileOps) Delete(fpath string) error {
}
