// Package experimental implements elsi.experimental module.
package experimental

import (
	"github.com/genkami/elsi/elrpc/helpers"
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
)

type Imports struct {
	Stream Stream
	File   File
	Stdio  Stdio
	HTTP   HTTP
}

type Exports struct{}

func UseWorld(instance types.Instance, imports *Imports) *Exports {
	ImportStream(instance, imports.Stream)
	ImportFile(instance, imports.File)
	ImportStdio(instance, imports.Stdio)
	ImportHTTP(instance, imports.HTTP)
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
	Close(handle *Handle) (*message.Void, error)
}

func ImportStream(instance types.Instance, stream Stream) {
	instance.Use(ModuleID, MethodID_Stream_Read, helpers.TypedHandler2[*Handle, *message.Uint64, *message.Bytes](stream.Read))
	instance.Use(ModuleID, MethodID_Stream_Write, helpers.TypedHandler2[*Handle, *message.Bytes, *message.Uint64](stream.Write))
	instance.Use(ModuleID, MethodID_Stream_Close, helpers.TypedHandler1[*Handle, *message.Void](stream.Close))
}

const (
	OpenModeCreate = 0
	OpenModeRead   = 1
	OpenModeWrite  = 2
	OpenModeAppend = 4
)

type File interface {
	Open(path *message.String, mode *message.Uint64) (*Handle, error)
}

func ImportFile(instance types.Instance, file File) {
	instance.Use(ModuleID, MethodID_File_Open, helpers.TypedHandler2[*message.String, *message.Uint64, *Handle](file.Open))
}

const (
	StdinID  = 0x0000_0000
	StdoutID = 0x0000_0001
	StderrID = 0x0000_0002
)

type Stdio interface {
	OpenStdHandle(id *message.Uint64) (*Handle, error)
}

func ImportStdio(instance types.Instance, stdio Stdio) {
	instance.Use(ModuleID, MethodID_Stdio_OpenStdHandle, helpers.TypedHandler1[*message.Uint64, *Handle](stdio.OpenStdHandle))
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

func ImportHTTP(instance types.Instance, http HTTP) {
	instance.Use(ModuleID, MethodID_HTTP_Listen, helpers.TypedHandler1[*message.String, *HTTPListener](http.Listen))
	instance.Use(ModuleID, MethodID_HTTP_CloseListener, helpers.TypedHandler1[*HTTPListener, *message.Void](http.CloseListener))
	instance.Use(ModuleID, MethodID_HTTP_PollRequest, helpers.TypedHandler1[*HTTPListener, *ServerRequest](http.PollRequest))
	instance.Use(ModuleID, MethodID_HTTP_SendResponse, helpers.TypedHandler2[*message.Uint64, *ServerResponseHeader, *Handle](http.SendResponseHeader))
}