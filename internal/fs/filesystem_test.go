package fs

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFS struct{}

func (t testFS) Open(name string) (fs.File, error) {
	return nil, os.ErrNotExist
}

func (t testFS) CloseFile(f fs.File) {
	// no-op
}

func TestOsFileSystem(t *testing.T) {
	fp, err := FS().Open(`../..//fixtures/.env`)
	assert.NoError(t, err)
	assert.NotNil(t, fp)

	FS().CloseFile(fp)
}

func TestFileSystem_API(t *testing.T) {
	Set(new(testFS))

	fp, err := FS().Open(`../..//fixtures/.env`)
	assert.Error(t, err)
	assert.Nil(t, fp)

	FS().CloseFile(fp)

	Reset()
}
