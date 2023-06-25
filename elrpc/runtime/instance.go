package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

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
	handlers map[uint64]types.Handler
	mod      Module
	exporter *exporterImpl
	wg       sync.WaitGroup
}

var _ types.Instance = (*Instance)(nil)

func NewInstance(mod Module) *Instance {
	exporter := newExporter()
	instance := &Instance{
		handlers: make(map[uint64]types.Handler),
		mod:      mod,
		exporter: exporter,
	}
	_ = builtin.UseWorld(instance, exporter)
	return instance
}

func (instance *Instance) Use(moduleID, methodID uint32, h types.Handler) {
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
		rlenBuf := make([]byte, message.LengthSize)
		_, err = io.ReadFull(stream, rlenBuf)
		if err != nil {
			return err
		}
		length, err := message.DecodeLength(rlenBuf)
		if err != nil {
			return err
		}

		req := make([]byte, length)
		_, err = io.ReadFull(stream, req)
		if err != nil {
			return err
		}
		dec := message.NewDecoder(req)

		resp := instance.dispatchRequest(dec)
		enc := message.NewEncoder()
		err = resp.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		respBody := enc.Buffer()

		wlenBuf, err := message.AppendLength(nil, len(respBody))
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

func (instance *Instance) dispatchRequest(dec *message.Decoder) *message.Result[message.Message, *message.Error] {
	type Resp = message.Result[message.Message, *message.Error]
	mID, err := dec.DecodeUint64()
	if err != nil {
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				Code:    builtin.CodeInvalidRequest,
				Message: "failed to decode method ID",
			},
		}
	}
	handler, ok := instance.handlers[mID]
	if !ok {
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				Code:    builtin.CodeUnimplemented,
				Message: fmt.Sprintf("method %X is not implemented", mID),
			},
		}
	}
	resp, err := handler.HandleRequest(dec)
	if err != nil {
		var elrpcErr *message.Error
		if errors.As(err, &elrpcErr) {
			return &Resp{IsOk: false, Err: elrpcErr}
		}
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				Code:    builtin.CodeInternal,
				Message: err.Error(),
			},
		}
	}
	return &Resp{IsOk: true, Ok: resp}
}

func (instance *Instance) Call(moduleID, methodID uint32, args *message.Any) (*message.Any, error) {
	ch := instance.exporter.callAsync(&builtin.MethodCall{
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
	retVal *message.Result[*message.Any, *message.Error]
}

type exporterImpl struct {
	mu        sync.Mutex
	waiters   map[uint64]chan<- callResult
	callQueue []*builtin.MethodCall
	next      uint64
}

var _ builtin.Exporter = &exporterImpl{}

func newExporter() *exporterImpl {
	return &exporterImpl{
		waiters: make(map[uint64]chan<- callResult),
	}
}

func (e *exporterImpl) callAsync(call *builtin.MethodCall) <-chan callResult {
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

func (e *exporterImpl) PollMethodCall() (*builtin.MethodCall, error) {
	type Resp = message.Result[*builtin.MethodCall, *message.Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.callQueue) == 0 {
		return nil, &message.Error{
			Code:    builtin.CodeNotFound,
			Message: "no method call",
		}
	}
	call := e.callQueue[0]
	e.callQueue = e.callQueue[1:]
	return call, nil
}

func (e *exporterImpl) SendResult(m *builtin.MethodResult) (*message.Void, error) {
	type Resp = message.Result[*message.Void, *message.Error]
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.waiters[m.CallID]
	if !ok {
		return nil, &message.Error{
			Code:    builtin.CodeNotFound,
			Message: "no such method call",
		}
	}
	ch <- callResult{m.RetVal}
	return &message.Void{}, nil
}

func fullID(moduleID, methodID uint32) uint64 {
	return uint64(moduleID)<<32 | uint64(methodID)
}
