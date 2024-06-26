package message

import "fmt"

type Option[T Message] struct {
	IsSome bool
	Some   T
}

var _ Message = (*Option[Message])(nil)

func (o *Option[T]) UnmarshalELRPC(dec *Decoder) error {
	vtag, err := dec.DecodeVariantTag()
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
		err = enc.EncodeVariantTag(0)
		if err != nil {
			return err
		}
		err = o.Some.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		return nil
	} else {
		return enc.EncodeVariantTag(1)
	}
}

func (o *Option[T]) ZeroMessage() Message {
	return &Option[T]{}
}

type Result[T, U Message] struct {
	IsOk bool
	Ok   T
	Err  U
}

var _ Message = (*Result[Message, Message])(nil)

func (e *Result[T, U]) UnmarshalELRPC(dec *Decoder) error {
	vtag, err := dec.DecodeVariantTag()
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

func (e *Result[T, U]) MarshalELRPC(enc *Encoder) error {
	var err error
	if e.IsOk {
		err = enc.EncodeVariantTag(0)
		if err != nil {
			return err
		}
		err = e.Ok.MarshalELRPC(enc)
		if err != nil {
			return err
		}
		return nil
	} else {
		err = enc.EncodeVariantTag(1)
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

func (e *Result[T, U]) ZeroMessage() Message {
	return &Result[T, U]{}
}

type Error struct {
	ModuleID uint32
	Code     uint32
	Message  string
}

var _ Message = (*Error)(nil)

func (e *Error) Error() string {
	if e == nil {
		return "elrpc: error (nil)"
	}
	return fmt.Sprintf("elrpc: error (mod = %X, code = %X): %s", e.ModuleID, e.Code, e.Message)
}

func (e *Error) UnmarshalELRPC(dec *Decoder) error {
	modID, err := dec.DecodeUint32()
	if err != nil {
		return err
	}
	code, err := dec.DecodeUint32()
	if err != nil {
		return err
	}
	msg, err := dec.DecodeString()
	if err != nil {
		return err
	}
	e.ModuleID = modID
	e.Code = code
	e.Message = msg
	return nil
}

func (e *Error) MarshalELRPC(enc *Encoder) error {
	err := enc.EncodeUint32(e.ModuleID)
	if err != nil {
		return err
	}
	err = enc.EncodeUint32(e.Code)
	if err != nil {
		return err
	}
	err = enc.EncodeString(e.Message)
	if err != nil {
		return err
	}
	return nil
}

func (e *Error) ZeroMessage() Message {
	return &Error{}
}
