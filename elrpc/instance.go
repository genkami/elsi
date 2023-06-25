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

type World struct {
	imports map[string]Handler
	// TODO: exports
}

func NewWorld(imports map[string]Handler) *World {
	w := &World{
		imports: make(map[string]Handler),
	}
	// explicitly copy the map
	for name, h := range imports {
		w.imports[name] = h
	}
	return w
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
	handlers map[string]Handler
	mod      Module
	exporter *exporterImpl
	wg       sync.WaitGroup
}

func NewInstance(mod Module) *Instance {
	exporter := newExporter()
	instance := &Instance{
		handlers: make(map[string]Handler),
		mod:      mod,
		exporter: exporter,
	}
	instance.Use(newBuiltinWorld(exporter))
	return instance
}

func (instance *Instance) Use(w *World) {
	for name, h := range w.imports {
		instance.handlers[name] = h
	}
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
	methodName, err := dec.DecodeBytes()
	if err != nil {
		// TODO: message
		return &Resp{IsOk: false, Err: &Error{Code: CodeInvalidRequest}}
	}
	handler, ok := instance.handlers[string(methodName)]
	if !ok {
		return &Resp{IsOk: false, Err: &Error{Code: CodeUnimplemented}}
	}
	resp, err := handler.HandleRequest(dec)
	if err != nil {
		// TODO: message
		return &Resp{IsOk: false, Err: &Error{Code: CodeInternal}}
	}
	return &Resp{IsOk: true, Ok: resp}
}

// TODO: there can be an error
func (instance *Instance) Call(name []byte, args *Any) *Any {
	ch := instance.exporter.callAsync(&MethodCall{
		Name: name,
		Args: args,
	})
	r := <-ch
	return r.retVal
}

type callResult struct {
	retVal *Any
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
	call.ID = id
	e.callQueue = append(e.callQueue, call)
	e.waiters[id] = ch
	return ch
}

func (e *exporterImpl) PollMethodCall() (*MethodCall, error) {
	type Resp = Result[*MethodCall, *Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.callQueue) == 0 {
		// TODO: message
		return nil, errors.New("no method call")
	}
	call := e.callQueue[0]
	e.callQueue = e.callQueue[1:]
	return call, nil
}

func (e *exporterImpl) SendResult(m *MethodResult) (*Void, error) {
	type Resp = Result[*Void, *Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.waiters[m.ID]
	if !ok {
		// TODO: message
		return nil, errors.New("not found")
	}
	ch <- callResult{m.RetVal}
	return &Void{}, nil
}
