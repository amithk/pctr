package pctr

import "os"

// If the file doesn't exists, create one.
func OpenFile(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// File must be closed before MoveFile is called on it.
func MoveFile(oldPath, newPath string) error {
	return nil
}

func ReadFile(f *os.File, buf []byte) error {
	_, err := f.ReadAt(buf, 0)
	return err
}

func WriteFile(f *os.File, buf []byte) error {
	_, err := f.WriteAt(buf, 0)
	return err
}

// File must be closed before DeleteFile is called on it.
func DeleteFile(fpath string) error {
	return os.RemoveAll(fpath)
}

func CloseFile(f *os.File) error {
	return nil
}
