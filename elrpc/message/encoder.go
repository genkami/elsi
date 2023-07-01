package message

func AppendLength(buf []byte, length int) ([]byte, error) {
	if length < 0 {
		return nil, ErrTooLarge
	}
	return endian.AppendUint64(buf, uint64(length)), nil
}

type Encoder struct {
	buf []byte
}

func NewEncoder() *Encoder {
	return &Encoder{
		buf: make([]byte, 0, 128),
	}
}

func (e *Encoder) EncodeUint8(val uint8) error {
	e.buf = append(e.buf, TagUint8, val)
	return nil
}

func (e *Encoder) EncodeUint16(val uint16) error {
	e.buf = append(e.buf, TagUint16)
	e.buf = endian.AppendUint16(e.buf, val)
	return nil
}

func (e *Encoder) EncodeUint32(val uint32) error {
	e.buf = append(e.buf, TagUint32)
	e.buf = endian.AppendUint32(e.buf, val)
	return nil
}

func (e *Encoder) EncodeUint64(val uint64) error {
	e.buf = append(e.buf, TagUint64)
	e.buf = endian.AppendUint64(e.buf, val)
	return nil
}

func (e *Encoder) EncodeInt8(val int8) error {
	e.buf = append(e.buf, TagInt8, uint8(val))
	return nil
}

func (e *Encoder) EncodeInt16(val int16) error {
	e.buf = append(e.buf, TagInt16)
	e.buf = endian.AppendUint16(e.buf, uint16(val))
	return nil
}

func (e *Encoder) EncodeInt32(val int32) error {
	e.buf = append(e.buf, TagInt32)
	e.buf = endian.AppendUint32(e.buf, uint32(val))
	return nil
}

func (e *Encoder) EncodeInt64(val int64) error {
	e.buf = append(e.buf, TagInt64)
	e.buf = endian.AppendUint64(e.buf, uint64(val))
	return nil
}

func (e *Encoder) EncodeBytes(val []byte) error {
	e.buf = append(e.buf, TagBytes)
	e.buf = endian.AppendUint64(e.buf, uint64(len(val)))
	e.buf = append(e.buf, val...)
	return nil
}

func (e *Encoder) EncodeString(val string) error {
	return e.EncodeBytes([]byte(val))
}

func (e *Encoder) EncodeArrayLen(val uint64) error {
	e.buf = append(e.buf, TagArray)
	e.buf = endian.AppendUint64(e.buf, val)
	return nil
}

func (e *Encoder) EncodeVariant(val uint8) error {
	e.buf = append(e.buf, TagVariant, val)
	return nil
}

func (e *Encoder) EncodeAny(val *Any) error {
	e.buf = append(e.buf, TagAny)
	e.buf = endian.AppendUint64(e.buf, uint64(len(val.Raw)))
	e.buf = append(e.buf, val.Raw...)
	return nil
}

func (e *Encoder) Buffer() []byte {
	return e.buf
}
