package expimpl

import (
	"io"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
)

type Stream struct {
	hs *HandleSet
}

var _ exp.Stream = (*Stream)(nil)

func (s *Stream) Read(handle *exp.Handle, size *message.Uint64) (*message.Bytes, error) {
	instance, ok := s.hs.Get(handle.ID)
	if !ok {
		return nil, errNoSuchHandle
	}
	r, ok := instance.(io.Reader)
	if !ok {
		return nil, errUnsupported
	}
	buf := make([]byte, size.Value)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		// TODO: convert to ELRPC error
		return nil, err
	}
	return &message.Bytes{Value: buf}, nil
}

func (s *Stream) Write(handle *exp.Handle, data *message.Bytes) (*message.Uint64, error) {
	instance, ok := s.hs.Get(handle.ID)
	if !ok {
		return nil, errNoSuchHandle
	}
	w, ok := instance.(io.Writer)
	if !ok {
		return nil, errUnsupported
	}
	n, err := w.Write(data.Value)
	if err != nil {
		// TODO: convert to ELRPC error
		return nil, err
	}
	return &message.Uint64{Value: uint64(n)}, nil
}

func (s *Stream) Close(handle *exp.Handle) (message.Void, error) {
	instance, ok := s.hs.Remove(handle.ID)
	if !ok {
		return message.Void{}, errNoSuchHandle
	}
	c, ok := instance.(io.Closer)
	if !ok {
		// Every handler can be closed.
		return message.Void{}, nil
	}
	err := c.Close()
	if err != nil {
		// TODO: convert to ELRPC error
		return message.Void{}, err
	}
	return message.Void{}, nil
}
