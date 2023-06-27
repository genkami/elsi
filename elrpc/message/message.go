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
// Any:        0x04 XX XX XX XX XX XX XX XX YY .. YY
//             where
//               XX .. XX = 64-bit big-endian integer (length)
//               YY .. YY = variable-length byte array representing another message

const (
	TagUint8   = 0x01
	TagUint16  = 0x02
	TagUint32  = 0x03
	TagUint64  = 0x04
	TagInt8    = 0x05
	TagInt16   = 0x06
	TagInt64   = 0x07
	TagInt32   = 0x08
	TagBytes   = 0x09
	TagVariant = 0x0a
	TagAny     = 0x0b
)

var (
	ErrTooLarge        = errors.New("size too large")
	ErrInsufficientBuf = errors.New("insufficient buffer")
	ErrTypeMismatch    = errors.New("type mismatch")
)

var endian interface {
	binary.ByteOrder
	binary.AppendByteOrder
} = binary.BigEndian

const LengthSize = 8

type Marshaler interface {
	MarshalELRPC(*Encoder) error
}

type Unmarshaler interface {
	UnmarshalELRPC(*Decoder) error
}

type Message interface {
	Unmarshaler
	Marshaler
	ZeroMessage() Message
}

func NewMessage[T Message]() Message {
	var z T
	return z.ZeroMessage()
}
