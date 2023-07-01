package message_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/google/go-cmp/cmp"
)

func TestVoid_UnmarshalELRPC(t *testing.T) {
	buf := []byte{}
	dec := message.NewDecoder(buf)
	var v message.Void
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVoid_MarshalELRPC(t *testing.T) {
	var v message.Void
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestVoid_ZeroMessage(t *testing.T) {
	var v message.Void
	got := v.ZeroMessage()
	if _, ok := got.(message.Void); !ok {
		t.Errorf("want Void but got %T", got)
	}
}

func TestUint8_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x01, 0xef}
	dec := message.NewDecoder(buf)
	var v message.Uint8
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != 0xef {
		t.Errorf("want 0xef but got 0x%X", v.Value)
	}
}

func TestUint8_MarshalELRPC(t *testing.T) {
	v := message.Uint8{Value: 0xef}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0xef}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestUint8_ZeroMessage(t *testing.T) {
	var v message.Uint8
	got := v.ZeroMessage()
	if _, ok := got.(*message.Uint8); !ok {
		t.Errorf("want Uint8 but got %T", got)
	}
}

func TestUint16_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x02, 0xef, 0xcd}
	dec := message.NewDecoder(buf)
	var v message.Uint16
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != 0xefcd {
		t.Errorf("want 0xefcd but got 0x%X", v.Value)
	}
}

func TestUint16_MarshalELRPC(t *testing.T) {
	v := message.Uint16{Value: 0xefcd}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0xef, 0xcd}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestUint16_ZeroMessage(t *testing.T) {
	var v message.Uint16
	got := v.ZeroMessage()
	if _, ok := got.(*message.Uint16); !ok {
		t.Errorf("want Uint16 but got %T", got)
	}
}

func TestUint32_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x03, 0xef, 0xcd, 0xab, 0x12}
	dec := message.NewDecoder(buf)
	var v message.Uint32
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != 0xefcdab12 {
		t.Errorf("want 0xefcdab12 but got 0x%X", v.Value)
	}
}

func TestUint32_MarshalELRPC(t *testing.T) {
	v := message.Uint32{Value: 0xefcdab12}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x03, 0xef, 0xcd, 0xab, 0x12}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestUint32_ZeroMessage(t *testing.T) {
	var v message.Uint32
	got := v.ZeroMessage()
	if _, ok := got.(*message.Uint32); !ok {
		t.Errorf("want Uint32 but got %T", got)
	}
}

func TestUint64_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x04, 0xef, 0xcd, 0xab, 0x12, 0x34, 0x56, 0x78, 0x90}
	dec := message.NewDecoder(buf)
	var v message.Uint64
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != 0xefcdab12_34567890 {
		t.Errorf("want 0xefcdab1234567890 but got 0x%X", v.Value)
	}
}

func TestUint64_MarshalELRPC(t *testing.T) {
	v := message.Uint64{Value: 0xefcdab12_34567890}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x04, 0xef, 0xcd, 0xab, 0x12, 0x34, 0x56, 0x78, 0x90}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestUint64_ZeroMessage(t *testing.T) {
	var v message.Uint64
	got := v.ZeroMessage()
	if _, ok := got.(*message.Uint64); !ok {
		t.Errorf("want Uint64 but got %T", got)
	}
}

func TestInt8_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x05, 0xef}
	dec := message.NewDecoder(buf)
	var v message.Int8
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != -0x11 {
		t.Errorf("want -0x11 but got 0x%X", v.Value)
	}
}

func TestInt8_MarshalELRPC(t *testing.T) {
	v := message.Int8{Value: -0x11}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x05, 0xef}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestInt8_ZeroMessage(t *testing.T) {
	var v message.Int8
	got := v.ZeroMessage()
	if _, ok := got.(*message.Int8); !ok {
		t.Errorf("want Int8 but got %T", got)
	}
}

func TestInt16_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x06, 0xef, 0xcd}
	dec := message.NewDecoder(buf)
	var v message.Int16
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != -0x1033 {
		t.Errorf("want -0x1033 but got 0x%X", v.Value)
	}
}

func TestInt16_MarshalELRPC(t *testing.T) {
	v := message.Int16{Value: -0x1033}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x06, 0xef, 0xcd}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestInt16_ZeroMessage(t *testing.T) {
	var v message.Int16
	got := v.ZeroMessage()
	if _, ok := got.(*message.Int16); !ok {
		t.Errorf("want Int16 but got %T", got)
	}
}

func TestInt32_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x07, 0xef, 0xcd, 0xab, 0x12}
	dec := message.NewDecoder(buf)
	var v message.Int32
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != -0x103254ee {
		t.Errorf("want -0x1033cded but got 0x%X", v.Value)
	}
}

func TestInt32_MarshalELRPC(t *testing.T) {
	v := message.Int32{Value: -0x1033cded}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x07, 0xef, 0xcc, 0x32, 0x13}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestInt32_ZeroMessage(t *testing.T) {
	var v message.Int32
	got := v.ZeroMessage()
	if _, ok := got.(*message.Int32); !ok {
		t.Errorf("want Int32 but got %T", got)
	}
}

func TestInt64_UnmarshalELRPC(t *testing.T) {
	buf := []byte{0x08, 0xef, 0xcd, 0xab, 0x12, 0x34, 0x56, 0x78, 0x9a}
	dec := message.NewDecoder(buf)
	var v message.Int64
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	if v.Value != -0x103254ed_cba98766 {
		t.Errorf("want -0x1033cdedab123456 but got 0x%X", v.Value)
	}
}

func TestInt64_MarshalELRPC(t *testing.T) {
	v := message.Int64{Value: -0x103254ed_cba98766}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x08, 0xef, 0xcd, 0xab, 0x12, 0x34, 0x56, 0x78, 0x9a}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestInt64_ZeroMessage(t *testing.T) {
	var v message.Int64
	got := v.ZeroMessage()
	if _, ok := got.(*message.Int64); !ok {
		t.Errorf("want Int64 but got %T", got)
	}
}
