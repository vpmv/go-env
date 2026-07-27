package fs

import (
	"io/fs"
	"os"
)

// FileSystem - internal FS wrapper
type FileSystem interface {
	Open(name string) (fs.File, error)
	CloseFile(f fs.File)
}

type osFileSystem struct{}

func (osFileSystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (osFileSystem) CloseFile(f fs.File) {
	if f != nil {
		_ = f.Close()
	}
}
