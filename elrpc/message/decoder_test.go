package message_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/google/go-cmp/cmp"
)

func TestDecodeLength(t *testing.T) {
	buf := []byte{0x00, 0x00, 0x00, 0x00, 0xde, 0xad, 0xbe, 0xef}
	n, err := message.DecodeLength(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0xdeadbeef {
		t.Errorf("want 0xdeadbeef but got 0x%x", n)
	}
}

func TestDecodeLength_insufficientBuf(t *testing.T) {
	buf := []byte{0xde, 0xad, 0xbe, 0xef}
	_, err := message.DecodeLength(buf)
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeInt64(t *testing.T) {
	buf := []byte{
		0x07,                                           // type tag (int64)
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := message.NewDecoder(buf)
	got, err := dec.DecodeInt64()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x123456789abcdef1
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeInt64_insufficientBuf(t *testing.T) {
	buf := []byte{0x07, 0x12, 0x34, 0x56, 0x78}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeInt64()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeInt64_typeMismatch(t *testing.T) {
	buf := []byte{
		0xff,                                           // type tag (int64)
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeInt64()
	if err != message.ErrTypeMismatch {
		t.Errorf("want ErrTypeMismatch but got %s", err)
	}
}

func TestDecoder_DecodeUint64(t *testing.T) {
	buf := []byte{
		0x04,                                           // type tag (uint64)
		0xf2, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := message.NewDecoder(buf)
	got, err := dec.DecodeUint64()
	if err != nil {
		t.Fatal(err)
	}

	var want uint64 = 0xf23456789abcdef1
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeUint64_insufficientBuf(t *testing.T) {
	buf := []byte{0x04, 0xf2, 0x34, 0x56, 0x78}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeUint64()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeUint64_typeMismatch(t *testing.T) {
	buf := []byte{
		0xff,                                           // type tag (uint64)
		0xf2, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeUint64()
	if err != message.ErrTypeMismatch {
		t.Errorf("want ErrTypeMismatch but got %s", err)
	}
}

func TestDecoder_DecodeBytes(t *testing.T) {
	buf := []byte{
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length = 10
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // value = "Hello"
	}
	dec := message.NewDecoder(buf)
	got, err := dec.DecodeBytes()
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("Hello")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeBytes_insufficientBuf_length(t *testing.T) {
	buf := []byte{
		0x09,                         // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, // length
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeBytes()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeBytes_insufficientBuf_body(t *testing.T) {
	buf := []byte{
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length = 10
		0x48, 0x65, 0x6c, 0x6c, // value = "Hell"
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeBytes()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeBytes_typeMismatch(t *testing.T) {
	buf := []byte{
		0xff,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length = 10
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // value = "Hello"
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeBytes()
	if err != message.ErrTypeMismatch {
		t.Errorf("want ErrTypeMismatch but got %s", err)
	}
}

func TestDecoder_DecodeString(t *testing.T) {
	buf := []byte{
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length = 10
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // value = "Hello"
	}
	dec := message.NewDecoder(buf)
	got, err := dec.DecodeString()
	if err != nil {
		t.Fatal(err)
	}

	want := "Hello"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeVariant(t *testing.T) {
	buf := []byte{
		0x0a, // type tag (variant)
		0xab, // value
	}
	dec := message.NewDecoder(buf)
	got, err := dec.DecodeVariant()
	if err != nil {
		t.Fatal(err)
	}

	var want uint8 = 0xab
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeVariant_insufficientBuf(t *testing.T) {
	buf := []byte{0x0a}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeVariant()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeVariant_typeMismatch(t *testing.T) {
	buf := []byte{
		0xff, // type tag
		0xab, // value
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeVariant()
	if err != message.ErrTypeMismatch {
		t.Errorf("want ErrTypeMismatch but got %s", err)
	}
}

func TestDecoder_DecodeAny(t *testing.T) {
	buf := []byte{
		0x0b,                                           // type tag (any)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, // length = 9
		0x07,                                           // type tag (int64)
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, // value
	}

	dec := message.NewDecoder(buf)
	any, err := dec.DecodeAny()
	if err != nil {
		t.Fatal(err)
	}

	anyDec := message.NewDecoder(any.Raw)
	got, err := anyDec.DecodeInt64()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x1122334455667788
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeAny_insufficientBuf_length(t *testing.T) {
	buf := []byte{
		0x0b,                         // type tag (any)
		0x00, 0x00, 0x00, 0x00, 0x00, // length
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeAny()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}

func TestDecoder_DecodeAny_insufficientBuf_body(t *testing.T) {
	buf := []byte{
		0x0b,                                           // type tag (any)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, // length = 9
		0x07,                               // type tag (int64)
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, // value
	}
	dec := message.NewDecoder(buf)
	_, err := dec.DecodeAny()
	if err != message.ErrInsufficientBuf {
		t.Errorf("want ErrInsufficientBuf but got %s", err)
	}
}
