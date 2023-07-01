package runtime

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/runtime/internal/builtinimpl"
	"github.com/genkami/elsi/elrpc/types"
)

type Instance struct {
	logger   *slog.Logger
	handlers map[uint64]types.Handler // a map from full method ID to its handler
	mod      Module
	exporter *builtinimpl.Exporter
	wg       sync.WaitGroup
}

var _ types.Instance = (*Instance)(nil)

func NewInstance(logger *slog.Logger, mod Module) *Instance {
	exporter := builtinimpl.NewExporter(logger)
	instance := &Instance{
		logger:   logger,
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
			// TODO: stop module
			instance.logger.Error("worker error", slog.String("error", err.Error()))
		}
	}()
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
		if !resp.IsOk {
			instance.logger.Error("method error", slog.String("error", resp.Err.Error()))
		}
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
	modID, err := dec.DecodeUint32()
	if err != nil {
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				ModuleID: builtin.ModuleID,
				Code:     builtin.CodeInvalidRequest,
				Message:  "failed to decode module ID",
			},
		}
	}
	methodID, err := dec.DecodeUint32()
	if err != nil {
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				ModuleID: builtin.ModuleID,
				Code:     builtin.CodeInvalidRequest,
				Message:  "failed to decode method ID",
			},
		}
	}
	fullMethodID := fullID(modID, methodID)
	handler, ok := instance.handlers[fullMethodID]
	if !ok {
		return &Resp{
			IsOk: false,
			Err: &message.Error{
				ModuleID: builtin.ModuleID,
				Code:     builtin.CodeUnimplemented,
				Message:  fmt.Sprintf("method %X in module %X is not implemented", modID, methodID),
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
				ModuleID: builtin.ModuleID,
				Code:     builtin.CodeInternal,
				Message:  err.Error(),
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
