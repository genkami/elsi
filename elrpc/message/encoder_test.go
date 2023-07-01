package message_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/google/go-cmp/cmp"
)

func TestEncoder_EncodeUint8(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint8(0xef)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x01, // type tag (uint8)
		0xef, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeUnt16(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint16(0xefcd)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x02,       // type tag (uint16)
		0xef, 0xcd, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeUint32(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(0xefcdab12)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x03,                   // type tag (uint32)
		0xef, 0xcd, 0xab, 0x12, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeUint64(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint64(0xfa2b3c4d5e6f7a8b)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x04,                                           // type tag (uint64)
		0xfa, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeInt8(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeInt8(0x12)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x05, // type tag (int8)
		0x12, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeInt16(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeInt16(0x12ab)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x06,       // type tag (int16)
		0x12, 0xab, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeInt32(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeInt32(0x12ab34cd)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x07,                   // type tag (int32)
		0x12, 0xab, 0x34, 0xcd, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeInt64(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeInt64(0x1a2b3c4d5e6f7a8b)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x08,                                           // type tag (int64)
		0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeBytes(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeBytes([]byte("Konnichiwa"))
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, // length = 10
		0x4b, 0x6f, 0x6e, 0x6e, 0x69, 0x63, 0x68, 0x69, 0x77, 0x61, // value = "Konnichiwa"
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeString(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeString("Konnichiwa")
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, // length = 10
		0x4b, 0x6f, 0x6e, 0x6e, 0x69, 0x63, 0x68, 0x69, 0x77, 0x61, // value = "Konnichiwa"
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeArrayLen(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeArrayLen(0xabcdef0123456789)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x0a,                                           // type tag (array)
		0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, // length
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)

	}
}

func TestEncoder_EncodeVariant(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeVariant(0xef)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x0b, // type tag (variant)
		0xef, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeAny(t *testing.T) {
	anyEnc := message.NewEncoder()
	err := anyEnc.EncodeBytes([]byte("Yo"))
	if err != nil {
		t.Fatal(err)
	}

	enc := message.NewEncoder()
	err = enc.EncodeAny(&message.Any{Raw: anyEnc.Buffer()})
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x0c,                                           // type tag (any)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, // length = ?
		0x09,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // length =2
		0x59, 0x6f, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
