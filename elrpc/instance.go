package elrpc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Handler interface {
	// TODO: this is not necessary
	DecodeRequest(*Decoder) (Message, error)
	HandleRequest(Message) Message
}

// TODO: make TypedFunc0, ..., TypedFunc5
type TypedHandlerFunc[Req, Resp Message] func(Req) Resp

func (h TypedHandlerFunc[Req, Resp]) DecodeRequest(dec *Decoder) (Message, error) {
	var z Req
	req := z.ZeroMessage()
	err := req.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h TypedHandlerFunc[Req, Resp]) HandleRequest(req Message) Message {
	return h(req.(Req))
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
	wg       sync.WaitGroup
}

func NewInstance(mod Module) *Instance {
	return &Instance{
		handlers: make(map[string]Handler),
		mod:      mod,
	}
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

		resp, err := instance.dispatchRequest(dec)
		if err != nil {
			return err
		}

		wlenBuf, err := AppendLength(nil, len(resp))
		if err != nil {
			return err
		}
		_, err = stream.Write(wlenBuf)
		if err != nil {
			return err
		}

		_, err = stream.Write(resp)
		if err != nil {
			return err
		}
	}
}

func (instance *Instance) dispatchRequest(dec *Decoder) ([]byte, error) {
	methodName, err := dec.DecodeBytes()
	if err != nil {
		return nil, err
	}
	handler, ok := instance.handlers[string(methodName)]
	if !ok {
		return nil, fmt.Errorf("no such method: %s", string(methodName))
	}
	req, err := handler.DecodeRequest(dec)
	if err != nil {
		return nil, err
	}
	resp := handler.HandleRequest(req)
	enc := NewEncoder()
	err = resp.MarshalELRPC(enc)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}
