package message

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

type Uint8 struct {
	Value uint8
}

var _ Message = (*Uint8)(nil)

func (u *Uint8) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeUint8()
	if err != nil {
		return err
	}
	u.Value = val
	return nil
}

func (u *Uint8) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeUint8(u.Value)
}

func (u *Uint8) ZeroMessage() Message {
	return &Uint8{}
}

type Uint16 struct {
	Value uint16
}

var _ Message = (*Uint16)(nil)

func (u *Uint16) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeUint16()
	if err != nil {
		return err
	}
	u.Value = val
	return nil
}

func (u *Uint16) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeUint16(u.Value)
}

func (u *Uint16) ZeroMessage() Message {
	return &Uint16{}
}

type Uint32 struct {
	Value uint32
}

var _ Message = (*Uint32)(nil)

func (u *Uint32) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeUint32()
	if err != nil {
		return err
	}
	u.Value = val
	return nil
}

func (u *Uint32) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeUint32(u.Value)
}

func (u *Uint32) ZeroMessage() Message {
	return &Uint32{}
}

type Uint64 struct {
	Value uint64
}

var _ Message = (*Uint64)(nil)

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

type Int8 struct {
	Value int8
}

var _ Message = (*Int8)(nil)

func (i *Int8) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeInt8()
	if err != nil {
		return err
	}
	i.Value = val
	return nil
}

func (i *Int8) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeInt8(i.Value)
}

func (i *Int8) ZeroMessage() Message {
	return &Int8{}
}

type Int16 struct {
	Value int16
}

var _ Message = (*Int16)(nil)

func (i *Int16) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeInt16()
	if err != nil {
		return err
	}
	i.Value = val
	return nil
}

func (i *Int16) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeInt16(i.Value)
}

func (i *Int16) ZeroMessage() Message {
	return &Int16{}
}

type Int32 struct {
	Value int32
}

var _ Message = (*Int32)(nil)

func (i *Int32) UnmarshalELRPC(dec *Decoder) error {
	val, err := dec.DecodeInt32()
	if err != nil {
		return err
	}
	i.Value = val
	return nil
}

func (i *Int32) MarshalELRPC(enc *Encoder) error {
	return enc.EncodeInt32(i.Value)
}

func (i *Int32) ZeroMessage() Message {
	return &Int32{}
}

type Int64 struct {
	Value int64
}

var _ Message = (*Int64)(nil)

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

var _ Message = (*Bytes)(nil)

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

var _ Message = (*String)(nil)

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

type Array[T Message] struct {
	Items []T
}

var _ Message = (*Array[Message])(nil)

func (a *Array[T]) UnmarshalELRPC(dec *Decoder) error {
	length, err := dec.DecodeArrayLen()
	if err != nil {
		return err
	}
	items := make([]T, length)
	for i := uint64(0); i < length; i++ {
		item := NewMessage[T]()
		err = item.UnmarshalELRPC(dec)
		if err != nil {
			return err
		}
		items[i] = item.(T)
	}
	a.Items = items
	return nil
}

func (a *Array[T]) MarshalELRPC(enc *Encoder) error {
	err := enc.EncodeArrayLen(uint64(len(a.Items)))
	if err != nil {
		return err
	}
	for _, item := range a.Items {
		err = item.MarshalELRPC(enc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Array[T]) ZeroMessage() Message {
	return &Array[T]{}
}

type Any struct {
	Raw []byte
}

var _ Message = (*Any)(nil)

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
