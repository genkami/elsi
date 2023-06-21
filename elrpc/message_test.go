package elrpc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/genkami/elsi/elrpc"
)

func TestDecoder_DecodeInt64(t *testing.T) {
	buf := []byte{
		0x00,                                           // type tag (int64)
		0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := elrpc.NewDecoder(buf)
	got, err := dec.DecodeInt64()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x123456789abcdef1
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeUint64(t *testing.T) {
	buf := []byte{
		0x02,                                           // type tag (uint64)
		0xf2, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf1, // value
	}
	dec := elrpc.NewDecoder(buf)
	got, err := dec.DecodeUint64()
	if err != nil {
		t.Fatal(err)
	}

	var want uint64 = 0xf23456789abcdef1
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeBytes(t *testing.T) {
	buf := []byte{
		0x01,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length = 10
		0x48, 0x65, 0x6c, 0x6c, 0x6f, // value = "Hello"
	}
	dec := elrpc.NewDecoder(buf)
	got, err := dec.DecodeBytes()
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("Hello")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestDecoder_DecodeVariant(t *testing.T) {
	buf := []byte{
		0x03, // type tag (variant)
		0xab, // value
	}
	dec := elrpc.NewDecoder(buf)
	got, err := dec.DecodeVariant()
	if err != nil {
		t.Fatal(err)
	}

	var want uint8 = 0xab
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeInt64(t *testing.T) {
	enc := elrpc.NewEncoder()
	err := enc.EncodeInt64(0x1a2b3c4d5e6f7a8b)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x00,                                           // type tag (int64)
		0x1a, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeUint64(t *testing.T) {
	enc := elrpc.NewEncoder()
	err := enc.EncodeUint64(0xfa2b3c4d5e6f7a8b)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x02,                                           // type tag (uint64)
		0xfa, 0x2b, 0x3c, 0x4d, 0x5e, 0x6f, 0x7a, 0x8b, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeBytes(t *testing.T) {
	enc := elrpc.NewEncoder()
	err := enc.EncodeBytes([]byte("Konnichiwa"))
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x01,                                           // type tag (bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, // length = 10
		0x4b, 0x6f, 0x6e, 0x6e, 0x69, 0x63, 0x68, 0x69, 0x77, 0x61, // value = "Konnichiwa"
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestEncoder_EncodeVariant(t *testing.T) {
	enc := elrpc.NewEncoder()
	err := enc.EncodeVariant(0xef)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0x03, // type tag (variant)
		0xef, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
