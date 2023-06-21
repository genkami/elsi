package elrpc

import (
	"errors"
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

type AnyTuple []Message

func (t AnyTuple) UnmarshalELRPC(dec *Decoder) error {
	// TODO: this sould be unmarshalable
	return errors.New("unmarshaling AnyTuple is not allowed")
}

func (t AnyTuple) MarshalELRPC(enc *Encoder) error {
	for _, m := range t {
		err := m.MarshalELRPC(enc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t AnyTuple) ZeroMessage() Message {
	return AnyTuple{}
}
