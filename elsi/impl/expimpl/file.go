package expimpl

import (
	"os"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
)

type File struct {
	hs *HandleSet
}

var _ exp.File = (*File)(nil)

func NewFile(hs *HandleSet) *File {
	return &File{hs: hs}
}

func (f *File) Open(path *message.String, mode *message.Uint64) (*exp.Handle, error) {
	openMode := buildOpenMode(mode.Value)
	// TODO: restrict access
	file, err := os.OpenFile(path.Value, openMode, 0644)
	if err != nil {
		// TODO: convert to ELRPC error
		return nil, err
	}
	hID := f.hs.Register(file)
	return &exp.Handle{ID: hID}, nil
}

func buildOpenMode(elsiMode uint64) int {
	var mode int
	if (elsiMode & exp.OpenModeCreate) != 0 {
		mode |= os.O_CREATE
	}
	if (elsiMode&exp.OpenModeRead) != 0 && (elsiMode&exp.OpenModeWrite) != 0 {
		mode |= os.O_RDWR
	} else if (elsiMode & exp.OpenModeRead) != 0 {
		mode |= os.O_RDONLY
	} else if (elsiMode & exp.OpenModeWrite) != 0 {
		mode |= os.O_WRONLY
	}
	if (elsiMode & exp.OpenModeAppend) != 0 {
		mode |= os.O_APPEND
	}
	return mode
}
