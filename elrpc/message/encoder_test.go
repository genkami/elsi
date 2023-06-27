package message_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/google/go-cmp/cmp"
)

func TestEncoder_EncodeInt64(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeInt64(0x1a2b3c4d5e6f7a8b)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x07,                                           // type tag (int64)
		0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b, // value
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

func TestEncoder_EncodeVariant(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeVariant(0xef)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x0a, // type tag (variant)
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
		0x0b,                                           // type tag (any)
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
