package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/runtime/internal/builtinimpl"
	"github.com/genkami/elsi/elrpc/types"
)

type Instance struct {
	// a map from full method ID to its handler
	handlers map[uint64]types.Handler
	mod      Module
	exporter *builtinimpl.Exporter
	wg       sync.WaitGroup
}

var _ types.Instance = (*Instance)(nil)

func NewInstance(mod Module) *Instance {
	exporter := builtinimpl.NewExporter()
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
	ch := instance.exporter.CallAsync(&builtin.MethodCall{
		ModuleID: moduleID,
		MethodID: methodID,
		Args:     args,
	})
	r := <-ch
	if !r.RetVal.IsOk {
		return nil, r.RetVal.Err
	}
	return r.RetVal.Ok, nil
}

func fullID(moduleID, methodID uint32) uint64 {
	return uint64(moduleID)<<32 | uint64(methodID)
}
