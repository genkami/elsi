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
	modID, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	methodID, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	args, err := dec.DecodeAny()
	if err != nil {
		return err
	}
	m.CallID = callID
	// TODO: uint32
	m.ModuleID = uint32(modID)
	m.MethodID = uint32(methodID)
	m.Args = args
	return nil
}

func (m *MethodCall) MarshalELRPC(enc *message.Encoder) error {
	err := enc.EncodeUint64(m.CallID)
	if err != nil {
		return err
	}
	// TODO: uint32
	err = enc.EncodeUint64(uint64(m.ModuleID))
	if err != nil {
		return err
	}
	err = enc.EncodeUint64(uint64(m.MethodID))
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
	SendResult(*MethodResult) (*message.Void, error)
}

func ImportExporter(instance types.Instance, e Exporter) {
	instance.Use(ModuleID, MethodID_Exporter_PollMethodCall, helpers.TypedHandler0[*MethodCall](e.PollMethodCall))
	instance.Use(ModuleID, MethodID_Exporter_SendResult, helpers.TypedHandler1[*MethodResult, *message.Void](e.SendResult))
}

type Exports struct{}

func UseWorld(instance types.Instance, e Exporter) *Exports {
	ImportExporter(instance, e)
	return &Exports{}
}
