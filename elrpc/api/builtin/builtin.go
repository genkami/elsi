package builtin

import (
	"github.com/genkami/elsi/elrpc/helpers"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

const (
	ModuleID = 0x0000_0000

	MethodID_Exporter_PollMethodCall = 0x0000_0000
	MethodID_Exporter_SendResult     = 0x0000_0001
)

const (
	CodeUnknown        = 0x0000
	CodeUnimplemented  = 0x0001
	CodeNotFound       = 0x0002
	CodeInvalidRequest = 0x0003
	CodeInternal       = 0x0004
)

type MethodCall struct {
	CallID   uint64
	ModuleID uint32
	MethodID uint32
	Args     *message.Any
}

func (m *MethodCall) UnmarshalELRPC(dec *message.Decoder) error {
	callID, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	modID, err := dec.DecodeUint32()
	if err != nil {
		return err
	}
	methodID, err := dec.DecodeUint32()
	if err != nil {
		return err
	}
	args, err := dec.DecodeAny()
	if err != nil {
		return err
	}
	m.CallID = callID
	m.ModuleID = modID
	m.MethodID = methodID
	m.Args = args
	return nil
}

func (m *MethodCall) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(m.CallID)
	if err != nil {
		return err
	}
	err = enc.EncodeUint32(m.ModuleID)
	if err != nil {
		return err
	}
	err = enc.EncodeUint32(m.MethodID)
	if err != nil {
		return err
	}
	err = enc.EncodeAny(m.Args)
	if err != nil {
		return err
	}
	return nil
}

func (m *MethodCall) ZeroMessage() message.Message {
	return &MethodCall{}
}

type MethodResult struct {
	CallID uint64
	RetVal *message.Result[*message.Any, *message.Error]
}

func (m *MethodResult) UnmarshalELRPC(dec *message.Decoder) error {
	id, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	resp := &message.Result[*message.Any, *message.Error]{}
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return err
	}
	m.CallID = id
	m.RetVal = resp
	return nil
}

func (m *MethodResult) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(m.CallID)
	if err != nil {
		return err
	}
	err = m.RetVal.MarshalELRPC(enc)
	if err != nil {
		return err
	}
	return nil
}

func (m *MethodResult) ZeroMessage() message.Message {
	return &MethodResult{}
}

type Exporter interface {
	PollMethodCall() (*MethodCall, error)
	SendResult(*MethodResult) (message.Void, error)
}

func ImportExporter(rt types.Runtime, e Exporter) {
	rt.Use(ModuleID, MethodID_Exporter_PollMethodCall, helpers.TypedHandler0[*MethodCall](e.PollMethodCall))
	rt.Use(ModuleID, MethodID_Exporter_SendResult, helpers.TypedHandler1[*MethodResult, message.Void](e.SendResult))
}

type Exports struct{}

func UseWorld(rt types.Runtime, e Exporter) *Exports {
	ImportExporter(rt, e)
	return &Exports{}
}
