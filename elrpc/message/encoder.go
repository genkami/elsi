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

func (e *Encoder) EncodeInt64(val int64) error {
	e.buf = append(e.buf, TagInt64)
	e.buf = endian.AppendUint64(e.buf, uint64(val))
	return nil
}

func (e *Encoder) EncodeUint64(val uint64) error {
	e.buf = append(e.buf, TagUint64)
	e.buf = endian.AppendUint64(e.buf, val)
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
