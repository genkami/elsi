package message

import (
	"encoding/binary"
	"errors"
)

// Types
//
// Int64:      0x00 XX XX XX XX XX XX XX XX
//             where
//               XX .. XX = 64-bit big-endian integer
// Bytes:      0x01 XX XX XX XX XX XX XX XX YY .. YY
//             where
//               XX .. XX = 64-bit big-endian integer (length)
//               YY .. YY = variable-length byte array
// Uint64:     0x03 XX XX XX XX XX XX XX XX
//             where
//               XX .. XX = 64-bit big-endian integer
// Variant:    0x04 XX YY .. YY
//             where
//               XX = 8-bit unsigned integer
//               YY .. YY = variable-length ELRPC message object

const (
	TagInt64   = 0x00
	TagBytes   = 0x01
	TagUint64  = 0x02
	TagVariant = 0x03
)

var (
	ErrTooLarge        = errors.New("size too large")
	ErrInsufficientBuf = errors.New("insufficient buffer")
	ErrTypeMismatch    = errors.New("type mismatch")
)

var endian = binary.BigEndian

const LengthSize = 8

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

func AppendLength(buf []byte, length int) ([]byte, error) {
	if length < 0 {
		return nil, ErrTooLarge
	}
	return endian.AppendUint64(buf, uint64(length)), nil
}

type Decoder struct {
	buf []byte
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{
		buf: buf,
	}
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

func (e *Encoder) EncodeVariant(val uint8) error {
	e.buf = append(e.buf, TagVariant, val)
	return nil
}

func (e *Encoder) Buffer() []byte {
	return e.buf
}
