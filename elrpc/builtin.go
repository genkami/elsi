package elrpc

import (
	"fmt"
)

type Void struct{}

func (Void) UnmarshalELRPC(dec *Decoder) error {
	return nil
}

func (Void) MarshalELRPC(enc *Encoder) error {
	return nil
}

func (Void) ZeroMessage() Message {
	return Void{}
}

type Option[T Message] struct {
	IsSome bool
	Some   T
}

func (o *Option[T]) UnmarshalELRPC(dec *Decoder) error {
	vtag, err := dec.DecodeVariant()
	if err != nil {
		return err
	}
	switch vtag {
	case 0:
		var z T
		someVal := z.ZeroMessage()
		err = someVal.UnmarshalELRPC(dec)
		if err != nil {
			return err
		}
		o.IsSome = true
		o.Some = someVal.(T)
		return nil
	case 1:
		o.IsSome = false
		return nil
	default:
		return fmt.Errorf("option: invalid variant: %d", vtag)
	}
}

func (o *Option[T]) MarshalELRPC(enc *Encoder) error {
	var err error
	if o.IsSome {
		err = enc.EncodeVariant(0)
		if err != nil {
			return err
		}
		err = o.Some.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		return nil
	} else {
		return enc.EncodeVariant(1)
	}
}

type Either[T, U Message] struct {
	IsOk bool
	Ok   T
	Err  U
}

func (e *Either[T, U]) UnmarshalELRPC(dec *Decoder) error {
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

func (e *Either[T, U]) MarshalELRPC(enc *Encoder) error {
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

const (
	CodeUnknown       = 0x0000
	CodeUnimplemented = 0x0001
	CodeNotFound      = 0x0002
)

type Error struct {
	Code uint64
	// TODO: add message?
}

func (e *Error) UnmarshalELRPC(dec *Decoder) error {
	code, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	e.Code = code
	return nil
}

func (e *Error) MarshalELRPC(enc *Encoder) error {
	err := enc.EncodeUint64(e.Code)
	if err != nil {
		return err
	}
	return nil
}

func (e *Error) ZeroMessage() Message {
	return &Error{}
}

type Any struct {
	Raw []byte
}

func (a *Any) UnmarshalELRPC(dec *Decoder) error {
	b, err := dec.DecodeAny()
	if err != nil {
		return err
	}
	a.Raw = b.Raw
	return nil
}

func (a *Any) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeAny(a)
}

func (a *Any) ZeroMessage() Message {
	return &Any{}
}

type MethodCall struct {
	ID   uint64
	Name []byte
	Args *Any
}

func (m *MethodCall) UnmarshalELRPC(dec *Decoder) error {
	id, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	name, err := dec.DecodeBytes()
	if err != nil {
		return err
	}
	args, err := dec.DecodeAny()
	if err != nil {
		return err
	}
	m.ID = id
	m.Name = name
	m.Args = args
	return nil
}

func (m *MethodCall) MarshalELRPC(enc *Encoder) error {
	err := enc.EncodeUint64(m.ID)
	if err != nil {
		return err
	}
	err = enc.EncodeBytes(m.Name)
	if err != nil {
		return err
	}
	err = enc.EncodeAny(m.Args)
	if err != nil {
		return err
	}
	return nil
}

func (m *MethodCall) ZeroMessage() Message {
	return &MethodCall{}
}

type MethodResult struct {
	ID     uint64
	RetVal *Any
}

func (m *MethodResult) UnmarshalELRPC(dec *Decoder) error {
	id, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	ret, err := dec.DecodeAny()
	if err != nil {
		return err
	}
	m.ID = id
	m.RetVal = ret
	return nil
}

func (m *MethodResult) MarshalELRPC(enc *Encoder) error {
	err := enc.EncodeUint64(m.ID)
	if err != nil {
		return err
	}
	err = enc.EncodeAny(m.RetVal)
	if err != nil {
		return err
	}
	return nil
}

func (m *MethodResult) ZeroMessage() Message {
	return &MethodResult{}
}

type Exporter interface {
	PollMethodCall() *Either[*MethodCall, *Error]
	SendResult(*MethodResult) *Either[*Void, *Error]
}

// The world that every ELRPC instance should use.
// This is automatically registered by system.
func newBuiltinWorld(e Exporter) *World {
	imports := map[string]Handler{
		"exporter/poll_method_call": TypedHandler0[*Either[*MethodCall, *Error]](e.PollMethodCall),
		"exporter/send_result":      TypedHandler1[*MethodResult, *Either[*Void, *Error]](e.SendResult),
	}
	return NewWorld(imports)
}
