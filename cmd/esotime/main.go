package main

import (
	"fmt"
	"os"

	"github.com/genkami/elsi/elrpc"
)

var theWorld *elrpc.World

func init() {
	w := elrpc.NewWorld()
	handlers := map[string]elrpc.AnyHandler{
		"elsi.x.ping":       &PingHandler,
		"elsi.x.add":        &AddHandler,
		"elsi.x.div":        &DivHandler,
		"elsi.x.write_file": &WriteFileHandler,
	}
	for name, h := range handlers {
		w.Register(name, h)
	}
	theWorld = w
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: esotime run CMD...\n")
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}
	if args[1] != "run" {
		usage()
	}

	mod := elrpc.NewProcessModule(args[2], args[3:]...)
	instance := elrpc.NewInstance(mod)
	instance.Use(theWorld)
	err := instance.Start()
	if err != nil {
		panic(err)
	}

	err = instance.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}

type Either[T, U elrpc.Message] struct {
	IsOk bool
	Ok   T
	Err  U
}

func (e *Either[T, U]) UnmarshalELRPC(dec *elrpc.Decoder) error {
	vtag, err := dec.DecodeVariant()
	if err != nil {
		return err
	}
	switch vtag {
	case 0:
		var z T
		okVal := z.ZeroMessage()
		err = okVal.UnmarshalELRPC(dec)
		if err != nil {
			return err
		}
		e.IsOk = true
		e.Ok = okVal.(T)
		return nil
	case 1:
		var z U
		errVal := z.ZeroMessage()
		err = errVal.UnmarshalELRPC(dec)
		if err != nil {
			return err
		}
		e.IsOk = false
		e.Err = errVal.(U)
		return nil
	default:
		return fmt.Errorf("either: invalid variant: %d", vtag)
	}
}

func (e *Either[T, U]) MarshalELRPC(enc *elrpc.Encoder) error {
	var err error
	if e.IsOk {
		err = enc.EncodeVariant(0)
		if err != nil {
			return err
		}
		err = e.Ok.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		return nil
	} else {
		err = enc.EncodeVariant(1)
		if err != nil {
			return err
		}
		err = e.Err.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		return nil
	}
}

func (e *Either[T, U]) ZeroMessage() elrpc.Message {
	return &Either[T, U]{}
}

type Error struct {
	Code uint64
	// TODO: add message?
}

func (e *Error) UnmarshalELRPC(dec *elrpc.Decoder) error {
	code, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	e.Code = code
	return nil
}

func (e *Error) MarshalELRPC(enc *elrpc.Encoder) error {
	err := enc.EncodeUint64(e.Code)
	if err != nil {
		return err
	}
	return nil
}

func (e *Error) ZeroMessage() elrpc.Message {
	return &Error{}
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

func pingImpl(req *PingRequest) *PingResponse {
	return &PingResponse{
		Nonce: req.Nonce,
	}
}

var PingHandler = elrpc.Handler[*PingRequest, *PingResponse]{
	Name: "elsi.x.ping",
	Impl: pingImpl,
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

func addImpl(req *AddRequest) *AddResponse {
	return &AddResponse{
		Sum: req.X + req.Y,
	}
}

var AddHandler = elrpc.Handler[*AddRequest, *AddResponse]{
	Name: "elsi.x.add",
	Impl: addImpl,
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

func divImpl(req *DivRequest) *Either[*DivResponse, *Error] {
	if req.Y == 0 {
		return &Either[*DivResponse, *Error]{
			IsOk: false,
			Err: &Error{
				Code: 0xababcdcd,
			},
		}
	}
	return &Either[*DivResponse, *Error]{
		IsOk: true,
		Ok: &DivResponse{
			Result: req.X / req.Y,
		},
	}
}

var DivHandler = elrpc.Handler[*DivRequest, *Either[*DivResponse, *Error]]{
	Name: "elsi.x.div",
	Impl: divImpl,
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
	var err error
	err = enc.EncodeUint64(uint64(r.Length))
	if err != nil {
		return err
	}
	return nil
}

func (r *WriteFileResponse) ZeroMessage() elrpc.Message {
	return &WriteFileResponse{}
}

func writeFileImpl(req *WriteFileRequest) *Either[*WriteFileResponse, *Error] {
	type Resp = Either[*WriteFileResponse, *Error]
	length, err := os.Stdout.Write(req.Buf)
	if err != nil {
		return &Resp{
			IsOk: false,
			Err: &Error{
				Code: 0x12345,
			},
		}
	}
	return &Resp{
		IsOk: true,
		Ok: &WriteFileResponse{
			Length: uint64(length),
		},
	}
}

var WriteFileHandler = elrpc.Handler[*WriteFileRequest, *Either[*WriteFileResponse, *Error]]{
	Name: "elsi.x.write_file",
	Impl: writeFileImpl,
}
