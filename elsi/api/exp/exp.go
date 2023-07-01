// Package exp implements elsi.exp module.
package exp

import (
	"github.com/genkami/elsi/elrpc/apibuilder"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

const (
	ModuleID = 0x0000_0001

	MethodID_Stream_Read         = 0x0000_0000
	MethodID_Stream_Write        = 0x0000_0001
	MethodID_Stream_Close        = 0x0000_0002
	MethodID_File_Open           = 0x0000_0010
	MethodID_Stdio_OpenStdHandle = 0x0000_0020
	MethodID_HTTP_Listen         = 0x0000_0030
	MethodID_HTTP_CloseListener  = 0x0000_0031
	MethodID_HTTP_PollRequest    = 0x0000_0032
	MethodID_HTTP_SendResponse   = 0x0000_0033

	// Attempted to do an unsupported operation on a handle (e.g. write to a read-only handle).
	CodeUnsupported = 0x0000_0001
)

type Imports struct {
	Stream Stream
	File   File
	Stdio  Stdio
	HTTP   HTTP
}

type Exports struct{}

func UseWorld(rt types.Runtime, imports *Imports) *Exports {
	ImportStream(rt, imports.Stream)
	ImportFile(rt, imports.File)
	ImportStdio(rt, imports.Stdio)
	ImportHTTP(rt, imports.HTTP)
	return &Exports{}
}

type Handle struct {
	ID uint64
}

func (h *Handle) UnmarshalELRPC(dec *message.Decoder) error {
	var err error
	h.ID, err = dec.DecodeUint64()
	if err != nil {
		return err
	}
	return nil
}

func (h *Handle) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(h.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handle) ZeroMessage() message.Message {
	return &Handle{}
}

type Stream interface {
	Read(handle *Handle, size *message.Uint64) (*message.Bytes, error)
	Write(handle *Handle, buf *message.Bytes) (*message.Uint64, error)
	Close(handle *Handle) (message.Void, error)
}

func ImportStream(rt types.Runtime, stream Stream) {
	rt.Use(ModuleID, MethodID_Stream_Read, apibuilder.HostHandler2[*Handle, *message.Uint64, *message.Bytes](stream.Read))
	rt.Use(ModuleID, MethodID_Stream_Write, apibuilder.HostHandler2[*Handle, *message.Bytes, *message.Uint64](stream.Write))
	rt.Use(ModuleID, MethodID_Stream_Close, apibuilder.HostHandler1[*Handle, message.Void](stream.Close))
}

const (
	OpenModeCreate = 1
	OpenModeRead   = 2
	OpenModeWrite  = 4
	OpenModeAppend = 3
)

type File interface {
	Open(path *message.String, mode *message.Uint64) (*Handle, error)
}

func ImportFile(rt types.Runtime, file File) {
	rt.Use(ModuleID, MethodID_File_Open, apibuilder.HostHandler2[*message.String, *message.Uint64, *Handle](file.Open))
}

const (
	HandleTypeStdin  = 0x00
	HandleTypeStdout = 0x01
	HandleTypeStderr = 0x02
)

type Stdio interface {
	OpenStdHandle(handleType *message.Uint8) (*Handle, error)
}

func ImportStdio(rt types.Runtime, stdio Stdio) {
	rt.Use(ModuleID, MethodID_Stdio_OpenStdHandle, apibuilder.HostHandler1[*message.Uint8, *Handle](stdio.OpenStdHandle))
}

type ServerRequest struct {
	RequestID uint64
	Method    string
	Path      string
	Body      *Handle
}

func (r *ServerRequest) UnmarshalELRPC(dec *message.Decoder) error {
	var err error
	r.RequestID, err = dec.DecodeUint64()
	if err != nil {
		return err
	}
	r.Method, err = dec.DecodeString()
	if err != nil {
		return err
	}
	r.Path, err = dec.DecodeString()
	if err != nil {
		return err
	}
	r.Body = new(Handle)
	err = r.Body.UnmarshalELRPC(dec)
	if err != nil {
		return err
	}
	return nil
}

func (r *ServerRequest) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(r.RequestID)
	if err != nil {
		return err
	}
	err = enc.EncodeString(r.Method)
	if err != nil {
		return err
	}
	err = enc.EncodeString(r.Path)
	if err != nil {
		return err
	}
	err = r.Body.MarshalELRPC(enc)
	if err != nil {
		return err
	}
	return nil
}

func (r *ServerRequest) ZeroMessage() message.Message {
	return &ServerRequest{}
}

type ServerResponseHeader struct {
	Status int64
}

func (r *ServerResponseHeader) UnmarshalELRPC(dec *message.Decoder) error {
	var err error
	r.Status, err = dec.DecodeInt64()
	if err != nil {
		return err
	}
	return nil
}

func (r *ServerResponseHeader) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeInt64(r.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *ServerResponseHeader) ZeroMessage() message.Message {
	return &ServerResponseHeader{}
}

type HTTPListener struct {
	ID uint64
}

func (s *HTTPListener) UnmarshalELRPC(dec *message.Decoder) error {
	var err error
	s.ID, err = dec.DecodeUint64()
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPListener) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(s.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPListener) ZeroMessage() message.Message {
	return &HTTPListener{}
}

type HTTP interface {
	Listen(addrAndPort *message.String) (*HTTPListener, error)
	CloseListener(listener *HTTPListener) (*message.Void, error)
	PollRequest(server *HTTPListener) (*ServerRequest, error)
	SendResponseHeader(reqID *message.Uint64, header *ServerResponseHeader) (*Handle, error)
}

func ImportHTTP(rt types.Runtime, http HTTP) {
	rt.Use(ModuleID, MethodID_HTTP_Listen, apibuilder.HostHandler1[*message.String, *HTTPListener](http.Listen))
	rt.Use(ModuleID, MethodID_HTTP_CloseListener, apibuilder.HostHandler1[*HTTPListener, *message.Void](http.CloseListener))
	rt.Use(ModuleID, MethodID_HTTP_PollRequest, apibuilder.HostHandler1[*HTTPListener, *ServerRequest](http.PollRequest))
	rt.Use(ModuleID, MethodID_HTTP_SendResponse, apibuilder.HostHandler2[*message.Uint64, *ServerResponseHeader, *Handle](http.SendResponseHeader))
}
