package message

func DecodeLength(buf []byte) (int, error) {
	if len(buf) < LengthSize {
		return 0, ErrInsufficientBuf
	}
	length := int(endian.Uint64(buf))
	if length < 0 {
		return 0, ErrTooLarge
	}
	return length, nil
}

type Decoder struct {
	buf []byte
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{
		buf: buf,
	}
}

func (d *Decoder) DecodeUint8() (uint8, error) {
	if len(d.buf) < 2 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagUint8 {
		return 0, ErrTypeMismatch
	}
	val := d.buf[1]
	d.buf = d.buf[2:]
	return val, nil
}

func (d *Decoder) DecodeUint16() (uint16, error) {
	if len(d.buf) < 3 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagUint16 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint16(d.buf[1:])
	d.buf = d.buf[3:]
	return val, nil
}

func (d *Decoder) DecodeUint32() (uint32, error) {
	if len(d.buf) < 5 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagUint32 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint32(d.buf[1:])
	d.buf = d.buf[5:]
	return val, nil
}

func (d *Decoder) DecodeUint64() (uint64, error) {
	if len(d.buf) < 9 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagUint64 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint64(d.buf[1:])
	d.buf = d.buf[9:]
	return val, nil
}

func (d *Decoder) DecodeInt8() (int8, error) {
	if len(d.buf) < 2 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagInt8 {
		return 0, ErrTypeMismatch
	}
	val := d.buf[1]
	d.buf = d.buf[2:]
	return int8(val), nil
}

func (d *Decoder) DecodeInt16() (int16, error) {
	if len(d.buf) < 3 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagInt16 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint16(d.buf[1:])
	d.buf = d.buf[3:]
	return int16(val), nil
}
func (d *Decoder) DecodeInt32() (int32, error) {
	if len(d.buf) < 5 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagInt32 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint32(d.buf[1:])
	d.buf = d.buf[5:]
	return int32(val), nil
}

func (d *Decoder) DecodeInt64() (int64, error) {
	if len(d.buf) < 9 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagInt64 {
		return 0, ErrTypeMismatch
	}
	val := endian.Uint64(d.buf[1:])
	d.buf = d.buf[9:]
	return int64(val), nil
}

func (d *Decoder) DecodeBytes() ([]byte, error) {
	if len(d.buf) < 1 {
		return nil, ErrInsufficientBuf
	}
	if d.buf[0] != TagBytes {
		return nil, ErrTypeMismatch
	}
	length, err := DecodeLength(d.buf[1:])
	if err != nil {
		return nil, err
	}

	d.buf = d.buf[1+LengthSize:]
	if len(d.buf) < length {
		return nil, ErrInsufficientBuf
	}
	val := d.buf[:length]
	d.buf = d.buf[length:]
	return val, nil
}

func (d *Decoder) DecodeString() (string, error) {
	val, err := d.DecodeBytes()
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (d *Decoder) DecodeVariant() (uint8, error) {
	if len(d.buf) < 2 {
		return 0, ErrInsufficientBuf
	}
	if d.buf[0] != TagVariant {
		return 0, ErrTypeMismatch
	}
	val := d.buf[1]
	d.buf = d.buf[2:]
	return val, nil
}

func (d *Decoder) DecodeAny() (*Any, error) {
	if len(d.buf) < 1 {
		return nil, ErrInsufficientBuf
	}
	if d.buf[0] != TagAny {
		return nil, ErrTypeMismatch
	}
	length, err := DecodeLength(d.buf[1:])
	if err != nil {
		return nil, err
	}

	d.buf = d.buf[1+LengthSize:]
	if len(d.buf) < length {
		return nil, ErrInsufficientBuf
	}
	val := d.buf[:length]
	d.buf = d.buf[length:]
	return &Any{Raw: val}, nil
}
