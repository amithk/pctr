package pctr

import "encoding/binary"
import "io"
import "os"

func CreateNew(fpath string) (*os.File, error) {
	buf := serializeUint64(0)
	f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	_, err1 := f.Write(buf)
	if err1 != nil {
		return nil, err
	}
	return f, nil
}

// File must be closed before MoveFile is called on it.
func MoveFile(oldPath, newPath string) error {
	return nil
}

func ReadFile(f *os.File) (uint64, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	_, err := f.Read(buf)
	if err == io.EOF {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return deserializeUint64(buf), nil
}

func WriteFile(f *os.File, val uint64) error {
	buf := serializeUint64(val)
	_, err := f.WriteAt(buf, 0)
	return err
}

// File must be closed before DeleteFile is called on it.
func DeleteFile(fpath string) error {
	return nil
}

func CloseFile(f *os.File) error {
	return nil
}

// TODO: Not sure if OpenFile is needed
func OpenFile(fpath string) (*os.File, error) {
	return nil, nil
}

func serializeUint64(val uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, 0)
	return buf
}

func deserializeUint64(buf []byte) uint64 {
	val, _ := binary.Uvarint(buf)
	return val
}
