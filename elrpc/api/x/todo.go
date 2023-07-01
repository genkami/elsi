package x

import (
	"github.com/genkami/elsi/elrpc/apibuilder"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

// TODO is an experimental interface that should be removed.
type TODO interface {
	Ping(*PingRequest) (*PingResponse, error)
	Add(*AddRequest) (*AddResponse, error)
	Div(*DivRequest) (*DivResponse, error)
	WriteFile(*WriteFileRequest) (*WriteFileResponse, error)
	TestExport() (*message.Void, error)
}

func ImportTODO(rt types.Runtime, todo TODO) {
	rt.Use(ModuleID, MethodID_TODO_Ping, apibuilder.HostHandler1[*PingRequest, *PingResponse](todo.Ping))
	rt.Use(ModuleID, MethodID_TODO_Add, apibuilder.HostHandler1[*AddRequest, *AddResponse](todo.Add))
	rt.Use(ModuleID, MethodID_TODO_Div, apibuilder.HostHandler1[*DivRequest, *DivResponse](todo.Div))
	rt.Use(ModuleID, MethodID_TODO_WriteFile, apibuilder.HostHandler1[*WriteFileRequest, *WriteFileResponse](todo.WriteFile))
	rt.Use(ModuleID, MethodID_TODO_TestExport, apibuilder.HostHandler0[*message.Void](todo.TestExport))
}

type PingRequest struct {
	Nonce int64
}

func (r *PingRequest) UnmarshalELRPC(dec *message.Decoder) error {
	nonce, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	r.Nonce = nonce
	return nil
}

func (r *PingRequest) MarshalELRPC(enc *message.Encoder) error {
	panic("PingRequest.MarshalELRPC: TODO")
}

func (r *PingRequest) ZeroMessage() message.Message {
	return &PingRequest{}
}

type PingResponse struct {
	Nonce int64
}

func (r *PingResponse) UnmarshalELRPC(dec *message.Decoder) error {
	panic("PingResponse.UnmarshalELRPC: TODO")
}

func (r *PingResponse) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeInt64(r.Nonce)
	if err != nil {
		return err
	}
	return nil
}

func (r *PingResponse) ZeroMessage() message.Message {
	return &PingResponse{}
}

type AddRequest struct {
	X, Y int64
}

func (r *AddRequest) UnmarshalELRPC(dec *message.Decoder) error {
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

func (r *AddRequest) MarshalELRPC(enc *message.Encoder) error {
	panic("AddRequest.MarshalELRPC: TODO")
}

func (r *AddRequest) ZeroMessage() message.Message {
	return &AddRequest{}
}

type AddResponse struct {
	Sum int64
}

func (r *AddResponse) UnmarshalELRPC(dec *message.Decoder) error {
	panic("AddResponse.UnmarshalELRPC: TODO")
}

func (r *AddResponse) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeInt64(r.Sum)
	if err != nil {
		return err
	}
	return nil
}

func (r *AddResponse) ZeroMessage() message.Message {
	return &AddResponse{}
}

type DivRequest struct {
	X, Y int64
}

func (r *DivRequest) UnmarshalELRPC(dec *message.Decoder) error {
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

func (r *DivRequest) MarshalELRPC(enc *message.Encoder) error {
	panic("DivRequest.MarshalELRPC: TODO")
}

func (r *DivRequest) ZeroMessage() message.Message {
	return &DivRequest{}
}

type DivResponse struct {
	Result int64
}

func (r *DivResponse) UnmarshalELRPC(dec *message.Decoder) error {
	panic("DivResponse.UnmarshalELRPC: TODO")
}

func (r *DivResponse) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeInt64(r.Result)
	if err != nil {
		return err
	}
	return nil
}

func (r *DivResponse) ZeroMessage() message.Message {
	return &DivResponse{}
}

type WriteFileRequest struct {
	Handle uint64
	Buf    []byte
}

func (r *WriteFileRequest) UnmarshalELRPC(dec *message.Decoder) error {
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

func (r *WriteFileRequest) MarshalELRPC(enc *message.Encoder) error {
	panic("WriteFileRequest.MarshalELRPC: TODO")
}

func (r *WriteFileRequest) ZeroMessage() message.Message {
	return &WriteFileRequest{}
}

type WriteFileResponse struct {
	Length uint64
}

func (r *WriteFileResponse) UnmarshalELRPC(dec *message.Decoder) error {
	panic("WriteFileResponse.UnmarshalELRPC: TODO")
}

func (r *WriteFileResponse) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(uint64(r.Length))
	if err != nil {
		return err
	}
	return nil
}

func (r *WriteFileResponse) ZeroMessage() message.Message {
	return &WriteFileResponse{}
}
