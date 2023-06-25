package x

import (
	"github.com/genkami/elsi/elrpc"
)

// TODO is an experimental interface that should be removed.
type TODO interface {
	Ping(*PingRequest) (*PingResponse, error)
	Add(*AddRequest) (*AddResponse, error)
	Div(*DivRequest) (*DivResponse, error)
	WriteFile(*WriteFileRequest) (*WriteFileResponse, error)
}

func ImportTODO(instance *elrpc.Instance, todo TODO) {
	instance.Use(ModuleID, MethodID_TODO_Ping, elrpc.TypedHandler1[*PingRequest, *PingResponse](todo.Ping))
	instance.Use(ModuleID, MethodID_TODO_Add, elrpc.TypedHandler1[*AddRequest, *AddResponse](todo.Add))
	instance.Use(ModuleID, MethodID_TODO_Div, elrpc.TypedHandler1[*DivRequest, *DivResponse](todo.Div))
	instance.Use(ModuleID, MethodID_TODO_WriteFile, elrpc.TypedHandler1[*WriteFileRequest, *WriteFileResponse](todo.WriteFile))
}

type PingRequest struct {
	Nonce int64
}

func (r *PingRequest) UnmarshalELRPC(dec *elrpc.Decoder) error {
	nonce, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	r.Nonce = nonce
	return nil
}

func (r *PingRequest) MarshalELRPC(enc *elrpc.Encoder) error {
	panic("PingRequest.MarshalELRPC: TODO")
}

func (r *PingRequest) ZeroMessage() elrpc.Message {
	return &PingRequest{}
}

type PingResponse struct {
	Nonce int64
}

func (r *PingResponse) UnmarshalELRPC(dec *elrpc.Decoder) error {
	panic("PingResponse.UnmarshalELRPC: TODO")
}

func (r *PingResponse) MarshalELRPC(enc *elrpc.Encoder) error {
	err := enc.EncodeInt64(r.Nonce)
	if err != nil {
		return err
	}
	return nil
}

func (r *PingResponse) ZeroMessage() elrpc.Message {
	return &PingResponse{}
}

type AddRequest struct {
	X, Y int64
}

func (r *AddRequest) UnmarshalELRPC(dec *elrpc.Decoder) error {
	x, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	y, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	r.X = x
	r.Y = y
	return nil
}

func (r *AddRequest) MarshalELRPC(enc *elrpc.Encoder) error {
	panic("AddRequest.MarshalELRPC: TODO")
}

func (r *AddRequest) ZeroMessage() elrpc.Message {
	return &AddRequest{}
}

type AddResponse struct {
	Sum int64
}

func (r *AddResponse) UnmarshalELRPC(dec *elrpc.Decoder) error {
	panic("AddResponse.UnmarshalELRPC: TODO")
}

func (r *AddResponse) MarshalELRPC(enc *elrpc.Encoder) error {
	err := enc.EncodeInt64(r.Sum)
	if err != nil {
		return err
	}
	return nil
}

func (r *AddResponse) ZeroMessage() elrpc.Message {
	return &AddResponse{}
}

type DivRequest struct {
	X, Y int64
}

func (r *DivRequest) UnmarshalELRPC(dec *elrpc.Decoder) error {
	x, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	y, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	r.X = x
	r.Y = y
	return nil
}

func (r *DivRequest) MarshalELRPC(enc *elrpc.Encoder) error {
	panic("DivRequest.MarshalELRPC: TODO")
}

func (r *DivRequest) ZeroMessage() elrpc.Message {
	return &DivRequest{}
}

type DivResponse struct {
	Result int64
}

func (r *DivResponse) UnmarshalELRPC(dec *elrpc.Decoder) error {
	panic("DivResponse.UnmarshalELRPC: TODO")
}

func (r *DivResponse) MarshalELRPC(enc *elrpc.Encoder) error {
	err := enc.EncodeInt64(r.Result)
	if err != nil {
		return err
	}
	return nil
}

func (r *DivResponse) ZeroMessage() elrpc.Message {
	return &DivResponse{}
}

type WriteFileRequest struct {
	Handle uint64
	Buf    []byte
}

func (r *WriteFileRequest) UnmarshalELRPC(dec *elrpc.Decoder) error {
	handle, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	buf, err := dec.DecodeBytes()
	if err != nil {
		return err
	}
	r.Handle = handle
	r.Buf = buf
	return nil
}

func (r *WriteFileRequest) MarshalELRPC(enc *elrpc.Encoder) error {
	panic("WriteFileRequest.MarshalELRPC: TODO")
}

func (r *WriteFileRequest) ZeroMessage() elrpc.Message {
	return &WriteFileRequest{}
}

type WriteFileResponse struct {
	Length uint64
}

func (r *WriteFileResponse) UnmarshalELRPC(dec *elrpc.Decoder) error {
	panic("WriteFileResponse.UnmarshalELRPC: TODO")
}

func (r *WriteFileResponse) MarshalELRPC(enc *elrpc.Encoder) error {
	err := enc.EncodeUint64(uint64(r.Length))
	if err != nil {
		return err
	}
	return nil
}

func (r *WriteFileResponse) ZeroMessage() elrpc.Message {
	return &WriteFileResponse{}
}
