package elrpc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Handler interface {
	HandleRequest(*Decoder) (Message, error)
}

type TypedHandler0[R Message] func() (R, error)

func (h TypedHandler0[R]) HandleRequest(dec *Decoder) (Message, error) {
	return h()
}

type TypedHandler1[T1, R Message] func(T1) (R, error)

func (h TypedHandler1[T1, R]) HandleRequest(dec *Decoder) (Message, error) {
	x1 := NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1))
}

type TypedHandler2[T1, T2, R Message] func(T1, T2) (R, error)

func (h TypedHandler2[T1, T2, R]) HandleRequest(dec *Decoder) (Message, error) {
	x1 := NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2))
}

type TypedHandler3[T1, T2, T3, R Message] func(T1, T2, T3) (R, error)

func (h TypedHandler3[T1, T2, T3, R]) HandleRequest(dec *Decoder) (Message, error) {
	x1 := NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3))
}

type TypedHandler4[T1, T2, T3, T4, R Message] func(T1, T2, T3, T4) (R, error)

func (h TypedHandler4[T1, T2, T3, T4, R]) HandleRequest(dec *Decoder) (Message, error) {
	x1 := NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x4 := NewMessage[T4]()
	err = x4.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3), x4.(T4))
}

type TypedHandler5[T1, T2, T3, T4, T5, R Message] func(T1, T2, T3, T4, T5) (R, error)

func (h TypedHandler5[T1, T2, T3, T4, T5, R]) HandleRequest(dec *Decoder) (Message, error) {
	x1 := NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x4 := NewMessage[T4]()
	err = x4.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x5 := NewMessage[T5]()
	err = x5.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3), x4.(T4), x5.(T5))
}

type methodCaller struct {
	instance *Instance
	moduleID uint32
	methodID uint32
}

type MethodCaller0[R Message] struct {
	methodCaller
}

func NewMethodCaller0[R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller0[R] {
	return &MethodCaller0[R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller0[R]) Call() (R, error) {
	var zero R
	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type MethodCaller1[T1, R Message] struct {
	methodCaller
}

func NewMethodCaller1[T1, R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller1[T1, R] {
	return &MethodCaller1[T1, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller1[T1, R]) Call(x1 T1) (R, error) {
	var zero R
	enc := NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type MethodCaller2[T1, T2, R Message] struct {
	methodCaller
}

func NewMethodCaller2[T1, T2, R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller2[T1, T2, R] {
	return &MethodCaller2[T1, T2, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller2[T1, T2, R]) Call(x1 T1, x2 T2) (R, error) {
	var zero R
	enc := NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type MethodCaller3[T1, T2, T3, R Message] struct {
	methodCaller
}

func NewMethodCaller3[T1, T2, T3, R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller3[T1, T2, T3, R] {
	return &MethodCaller3[T1, T2, T3, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller3[T1, T2, T3, R]) Call(x1 T1, x2 T2, x3 T3) (R, error) {
	var zero R
	enc := NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type MethodCaller4[T1, T2, T3, T4, R Message] struct {
	methodCaller
}

func NewMethodCaller4[T1, T2, T3, T4, R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller4[T1, T2, T3, T4, R] {
	return &MethodCaller4[T1, T2, T3, T4, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller4[T1, T2, T3, T4, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4) (R, error) {
	var zero R
	enc := NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x4.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type MethodCaller5[T1, T2, T3, T4, T5, R Message] struct {
	methodCaller
}

func NewMethodCaller5[T1, T2, T3, T4, T5, R Message](instance *Instance, moduleID, methodID uint32) *MethodCaller5[T1, T2, T3, T4, T5, R] {
	return &MethodCaller5[T1, T2, T3, T4, T5, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller5[T1, T2, T3, T4, T5, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4, x5 T5) (R, error) {
	var zero R
	enc := NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x4.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x5.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := NewDecoder(rawResp.Raw)
	resp := NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type Stream interface {
	io.ReadWriter
}

type pipeStream struct {
	io.Reader
	io.Writer
}

func NewPipeStream(in io.Reader, out io.Writer) Stream {
	return &pipeStream{in, out}
}

type Module interface {
	Stream() Stream
	Start() error
	Wait() error
}

type ProcessModule struct {
	cmd    *exec.Cmd
	stream Stream
}

func NewProcessModule(name string, args ...string) *ProcessModule {
	cmd := exec.Command(name, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr
	return &ProcessModule{
		cmd:    cmd,
		stream: NewPipeStream(stdout, stdin),
	}
}

func (m *ProcessModule) Stream() Stream {
	return m.stream
}

func (m *ProcessModule) Start() error {
	return m.cmd.Start()
}

func (m *ProcessModule) Wait() error {
	return m.cmd.Wait()
}

type Instance struct {
	// a map from full method ID to its handler
	handlers map[uint64]Handler
	mod      Module
	exporter *exporterImpl
	wg       sync.WaitGroup
}

func NewInstance(mod Module) *Instance {
	exporter := newExporter()
	instance := &Instance{
		handlers: make(map[uint64]Handler),
		mod:      mod,
		exporter: exporter,
	}
	_ = UseWorld(instance, exporter)
	return instance
}

func (instance *Instance) Use(moduleID, methodID uint32, h Handler) {
	instance.handlers[fullID(moduleID, methodID)] = h
}

func (instance *Instance) Start() error {
	err := instance.mod.Start()
	if err != nil {
		return err
	}

	instance.wg.Add(1)
	go func() {
		defer instance.wg.Done()
		err := instance.serverWorker()
		if err != nil {
			// TODO
			fmt.Fprintf(os.Stderr, "worker error: %s\n", err.Error())
		}
	}()
	// TODO: negotiation required?
	return nil
}

func (instance *Instance) Wait() error {
	// TODO: any way to terminate instead of waiting?
	err := instance.mod.Wait()
	if err != nil {
		return err
	}
	instance.wg.Wait()
	return nil
}

func (instance *Instance) serverWorker() error {
	var err error
	stream := instance.mod.Stream()
	for {
		rlenBuf := make([]byte, LengthSize)
		_, err = io.ReadFull(stream, rlenBuf)
		if err != nil {
			return err
		}
		length, err := DecodeLength(rlenBuf)
		if err != nil {
			return err
		}

		req := make([]byte, length)
		_, err = io.ReadFull(stream, req)
		if err != nil {
			return err
		}
		dec := NewDecoder(req)

		resp := instance.dispatchRequest(dec)
		enc := NewEncoder()
		err = resp.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		respBody := enc.Buffer()

		wlenBuf, err := AppendLength(nil, len(respBody))
		if err != nil {
			return err
		}
		_, err = stream.Write(wlenBuf)
		if err != nil {
			return err
		}

		_, err = stream.Write(respBody)
		if err != nil {
			return err
		}
	}
}

func (instance *Instance) dispatchRequest(dec *Decoder) *Result[Message, *Error] {
	type Resp = Result[Message, *Error]
	mID, err := dec.DecodeUint64()
	if err != nil {
		return &Resp{
			IsOk: false,
			Err: &Error{
				Code:    CodeInvalidRequest,
				Message: "failed to decode method ID",
			},
		}
	}
	handler, ok := instance.handlers[mID]
	if !ok {
		return &Resp{
			IsOk: false,
			Err: &Error{
				Code:    CodeUnimplemented,
				Message: fmt.Sprintf("method %X is not implemented", mID),
			},
		}
	}
	resp, err := handler.HandleRequest(dec)
	if err != nil {
		var elrpcErr *Error
		if errors.As(err, &elrpcErr) {
			return &Resp{IsOk: false, Err: elrpcErr}
		}
		return &Resp{
			IsOk: false,
			Err: &Error{
				Code:    CodeInternal,
				Message: err.Error(),
			},
		}
	}
	return &Resp{IsOk: true, Ok: resp}
}

func (instance *Instance) Call(moduleID, methodID uint32, args *Any) (*Any, error) {
	ch := instance.exporter.callAsync(&MethodCall{
		FullMethodID: fullID(moduleID, methodID),
		Args:         args,
	})
	r := <-ch
	if !r.retVal.IsOk {
		return nil, r.retVal.Err
	}
	return r.retVal.Ok, nil
}

type callResult struct {
	retVal *Result[*Any, *Error]
}

type exporterImpl struct {
	mu        sync.Mutex
	waiters   map[uint64]chan<- callResult
	callQueue []*MethodCall
	next      uint64
}

var _ Exporter = &exporterImpl{}

func newExporter() *exporterImpl {
	return &exporterImpl{
		waiters: make(map[uint64]chan<- callResult),
	}
}

func (e *exporterImpl) callAsync(call *MethodCall) <-chan callResult {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := make(chan callResult, 1)
	id := e.next
	e.next++
	call.CallID = id
	e.callQueue = append(e.callQueue, call)
	e.waiters[id] = ch
	return ch
}

func (e *exporterImpl) PollMethodCall() (*MethodCall, error) {
	type Resp = Result[*MethodCall, *Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.callQueue) == 0 {
		return nil, &Error{
			Code:    CodeNotFound,
			Message: "no method call",
		}
	}
	call := e.callQueue[0]
	e.callQueue = e.callQueue[1:]
	return call, nil
}

func (e *exporterImpl) SendResult(m *MethodResult) (*Void, error) {
	type Resp = Result[*Void, *Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.waiters[m.CallID]
	if !ok {
		return nil, &Error{
			Code:    CodeNotFound,
			Message: "no such method call",
		}
	}
	ch <- callResult{m.RetVal}
	return &Void{}, nil
}

func fullID(moduleID, methodID uint32) uint64 {
	return uint64(moduleID)<<32 | uint64(methodID)
}
