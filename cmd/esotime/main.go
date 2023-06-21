package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/genkami/elsi/elrpc"
)

var methodsMap map[string]AnyHandler

func init() {
	methodsMap = map[string]AnyHandler{}
	handlers := []AnyHandler{
		&PingHandler,
		&AddHandler,
		&DivHandler,
		&WriteFileHandler,
	}
	for _, h := range handlers {
		methodsMap[h.MethodName()] = h
	}
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
	cmd := exec.Command(args[2], args[3:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		strm := &pipeStream{Writer: stdin, Reader: stdout}
		err := serverWorker(strm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worker error: %s\n", err.Error())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: waiting worker goroutine to finish...\n")
	wg.Wait()

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}

type Stream interface {
	io.ReadWriter
}

type pipeStream struct {
	io.Reader
	io.Writer
}

type AnyHandler interface {
	MethodName() string
	DecodeRequest(*elrpc.Decoder) (Message, error)
	HandleRequest(Message) Message
}

type Message interface {
	elrpc.Unmarshaler
	elrpc.Marshaler
	ZeroMessage() Message
}

type Handler[Req, Resp Message] struct {
	Name string
	Impl func(Req) Resp
}

func (h *Handler[Req, Resp]) MethodName() string {
	return h.Name
}

func (h *Handler[Req, Resp]) DecodeRequest(dec *elrpc.Decoder) (Message, error) {
	var z Req
	req := z.ZeroMessage()
	err := req.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h *Handler[Req, Resp]) HandleRequest(req Message) Message {
	return h.Impl(req.(Req))
}

func serverWorker(strm Stream) error {
	var err error
	for {
		rlenBuf := make([]byte, elrpc.LengthSize)
		_, err = io.ReadFull(strm, rlenBuf)
		if err != nil {
			return err
		}
		length, err := elrpc.DecodeLength(rlenBuf)
		if err != nil {
			return err
		}

		req := make([]byte, length)
		_, err = io.ReadFull(strm, req)
		if err != nil {
			return err
		}
		dec := elrpc.NewDecoder(req)

		resp, err := dispatchRequest(dec)
		if err != nil {
			return err
		}

		wlenBuf, err := elrpc.AppendLength(nil, len(resp))
		if err != nil {
			return err
		}
		_, err = strm.Write(wlenBuf)
		if err != nil {
			return err
		}

		_, err = strm.Write(resp)
		if err != nil {
			return err
		}
	}
}

func dispatchRequest(dec *elrpc.Decoder) ([]byte, error) {
	methodName, err := dec.DecodeBytes()
	if err != nil {
		return nil, err
	}
	handler, ok := methodsMap[string(methodName)]
	if !ok {
		return nil, fmt.Errorf("no such method: %s", string(methodName))
	}
	req, err := handler.DecodeRequest(dec)
	if err != nil {
		return nil, err
	}
	resp := handler.HandleRequest(req)
	enc := elrpc.NewEncoder()
	err = resp.MarshalELRPC(enc)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}

type Either[T, U Message] struct {
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

func (e *Either[T, U]) ZeroMessage() Message {
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

func (e *Error) ZeroMessage() Message {
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

func (r *PingRequest) ZeroMessage() Message {
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

func (r *PingResponse) ZeroMessage() Message {
	return &PingResponse{}
}

func pingImpl(req *PingRequest) *PingResponse {
	return &PingResponse{
		Nonce: req.Nonce,
	}
}

var PingHandler = Handler[*PingRequest, *PingResponse]{
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

func (r *AddRequest) ZeroMessage() Message {
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

func (r *AddResponse) ZeroMessage() Message {
	return &AddResponse{}
}

func addImpl(req *AddRequest) *AddResponse {
	return &AddResponse{
		Sum: req.X + req.Y,
	}
}

var AddHandler = Handler[*AddRequest, *AddResponse]{
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

func (r *DivRequest) ZeroMessage() Message {
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

func (r *DivResponse) ZeroMessage() Message {
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

var DivHandler = Handler[*DivRequest, *Either[*DivResponse, *Error]]{
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

func (r *WriteFileRequest) ZeroMessage() Message {
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

func (r *WriteFileResponse) ZeroMessage() Message {
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

var WriteFileHandler = Handler[*WriteFileRequest, *Either[*WriteFileResponse, *Error]]{
	Name: "elsi.x.write_file",
	Impl: writeFileImpl,
}
