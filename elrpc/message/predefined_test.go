package message_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/google/go-cmp/cmp"
)

func TestOption_UnmarshalELRPC_some(t *testing.T) {
	buf := []byte{0x0b, 0x00, 0x01, 0xef}
	dec := message.NewDecoder(buf)
	got := message.NewMessage[*message.Option[*message.Uint8]]()
	err := got.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := &message.Option[*message.Uint8]{
		IsSome: true,
		Some:   &message.Uint8{Value: 0xef},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestOption_UnmarshalELRPC_none(t *testing.T) {
	buf := []byte{0x0b, 0x01}
	dec := message.NewDecoder(buf)
	got := message.NewMessage[*message.Option[*message.Uint8]]()
	err := got.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := &message.Option[*message.Uint8]{
		IsSome: false,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestOption_MarshalELRPC_some(t *testing.T) {
	v := &message.Option[*message.Uint8]{
		IsSome: true,
		Some:   &message.Uint8{Value: 0xef},
	}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x0b, 0x00, 0x01, 0xef}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestOption_MarshalELRPC_none(t *testing.T) {
	v := &message.Option[*message.Uint8]{
		IsSome: false,
	}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x0b, 0x01}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestOption_ZeroMessage(t *testing.T) {
	v := message.NewMessage[*message.Option[*message.Uint8]]()
	got := v.ZeroMessage()
	if _, ok := got.(*message.Option[*message.Uint8]); !ok {
		t.Errorf("want Uint8 but got %T", got)
	}
}

func TestResult_UnmarshalELRPC_ok(t *testing.T) {
	type Result = message.Result[*message.Uint8, *message.Uint16]
	buf := []byte{0x0b, 0x00, 0x01, 0xef}
	dec := message.NewDecoder(buf)
	got := message.NewMessage[*Result]()
	err := got.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := &Result{
		IsOk: true,
		Ok:   &message.Uint8{Value: 0xef},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestResult_UnmarshalELRPC_err(t *testing.T) {
	type Result = message.Result[*message.Uint8, *message.Uint16]
	buf := []byte{0x0b, 0x01, 0x02, 0xcd, 0xef}
	dec := message.NewDecoder(buf)
	got := message.NewMessage[*Result]()
	err := got.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := &Result{
		IsOk: false,
		Err:  &message.Uint16{Value: 0xcdef},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestResult_MarshalELRPC_ok(t *testing.T) {
	type Result = message.Result[*message.Uint8, *message.Uint16]
	v := &Result{
		IsOk: true,
		Ok:   &message.Uint8{Value: 0xef},
	}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x0b, 0x00, 0x01, 0xef}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestResult_MarshalELRPC_err(t *testing.T) {
	type Result = message.Result[*message.Uint8, *message.Uint16]
	v := &Result{
		IsOk: false,
		Err:  &message.Uint16{Value: 0xcdef},
	}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x0b, 0x01, 0x02, 0xcd, 0xef}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestResult_ZeroMessage(t *testing.T) {
	type Result = message.Result[*message.Uint8, *message.Uint16]
	v := message.NewMessage[*Result]()
	got := v.ZeroMessage()
	if _, ok := got.(*Result); !ok {
		t.Errorf("want Result but got %T", got)
	}
}
