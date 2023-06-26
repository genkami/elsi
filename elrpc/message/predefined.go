package message

import "fmt"

type Uint64 struct {
	Value uint64
}

func (u *Uint64) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeUint64()
	if err != nil {
		return err
	}
	u.Value = val
	return nil
}

func (u *Uint64) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeUint64(u.Value)
}

func (u *Uint64) ZeroMessage() Message {
	return &Uint64{}
}

type Int64 struct {
	Value int64
}

func (i *Int64) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeInt64()
	if err != nil {
		return err
	}
	i.Value = val
	return nil
}

func (i *Int64) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeInt64(i.Value)
}

func (i *Int64) ZeroMessage() Message {
	return &Int64{}
}

type Bytes struct {
	Value []byte
}

func (b *Bytes) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeBytes()
	if err != nil {
		return err
	}
	b.Value = val
	return nil
}

func (b *Bytes) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeBytes(b.Value)
}

func (b *Bytes) ZeroMessage() Message {
	return &Bytes{}
}

type String struct {
	Value string
}

func (s *String) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeString()
	if err != nil {
		return err
	}
	s.Value = val
	return nil
}

func (s *String) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeString(s.Value)
}

func (s *String) ZeroMessage() Message {
	return &String{}
}

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

type Result[T, U Message] struct {
	IsOk bool
	Ok   T
	Err  U
}

func (e *Result[T, U]) UnmarshalELRPC(dec *Decoder) error {
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

func (e *Result[T, U]) MarshalELRPC(enc *Encoder) error {
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

func (e *Result[T, U]) ZeroMessage() Message {
	return &Result[T, U]{}
}

type Error struct {
	Code    uint64
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("elrpc: error (code = %X): %s", e.Code, e.Message)
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
