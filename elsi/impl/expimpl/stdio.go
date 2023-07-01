package expimpl

import (
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
)

type StdHandleCtor = func() (any, error)

type Stdio struct {
	hs         *HandleSet
	stdHandles map[uint8]StdHandleCtor
}

var _ exp.Stdio = (*Stdio)(nil)

func NewStdio(hs *HandleSet, stdHandles map[uint8]StdHandleCtor) *Stdio {
	handles := make(map[uint8]StdHandleCtor)
	for k, v := range stdHandles {
		handles[k] = v
	}
	return &Stdio{
		hs:         hs,
		stdHandles: handles,
	}
}

func (s *Stdio) OpenStdHandle(hType *message.Uint8) (*exp.Handle, error) {
	ctor, ok := s.stdHandles[hType.Value]
	if !ok {
		return nil, errInvalidHandleType
	}
	instance, err := ctor()
	if err != nil {
		return nil, err
	}
	hID := s.hs.Register(instance)
	return &exp.Handle{ID: hID}, nil
}
