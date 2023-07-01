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

func TestBytes_UnmarshalELRPC(t *testing.T) {
	buf := []byte{
		0x09,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length
		0x12, 0x34, 0x56, 0x78, 0x9a, // value
	}
	dec := message.NewDecoder(buf)
	var v message.Bytes
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x12, 0x34, 0x56, 0x78, 0x9a}
	got := v.Value
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestBytes_MarshalELRPC(t *testing.T) {
	v := message.Bytes{Value: []byte{0x12, 0x34, 0x56, 0x78, 0x9a}}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x09,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length
		0x12, 0x34, 0x56, 0x78, 0x9a, // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestBytes_ZeroMessage(t *testing.T) {
	var v message.Bytes
	got := v.ZeroMessage()
	if _, ok := got.(*message.Bytes); !ok {
		t.Errorf("want Bytes but got %T", got)
	}
}

func TestString_UnmarshalELRPC(t *testing.T) {
	buf := []byte{
		0x09,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length
		'h', 'e', 'l', 'l', 'o', // value
	}
	dec := message.NewDecoder(buf)
	var v message.String
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := "hello"
	got := v.Value
	if want != got {
		t.Errorf("want %q but got %q", want, got)
	}
}

func TestString_MarshalELRPC(t *testing.T) {
	v := message.String{Value: "hello"}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x09,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, // length
		'h', 'e', 'l', 'l', 'o', // value
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestString_ZeroMessage(t *testing.T) {
	var v message.String
	got := v.ZeroMessage()
	if _, ok := got.(*message.String); !ok {
		t.Errorf("want String but got %T", got)
	}
}

func TestArray_UnmarshalELRPC(t *testing.T) {
	buf := []byte{
		0x0a,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // length
		0x01, 0xaa, // [0]: uint8 0xaa
		0x01, 0xbb, // [1]: uint8 0xbb
	}
	dec := message.NewDecoder(buf)
	var v message.Array[*message.Uint8]
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := []*message.Uint8{
		{Value: 0xaa},
		{Value: 0xbb},
	}
	got := v.Items
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestArray_UnmarshalELRPC_typeMismatch(t *testing.T) {
	buf := []byte{
		0x0a,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // length
		0x01, 0xaa, // [0]: uint8 0xaa
		0x02, 0xbb, 0xbb, // [1]: uint16 0xbbbb (type mismatch)
	}
	dec := message.NewDecoder(buf)
	var v message.Array[*message.Uint8]
	err := v.UnmarshalELRPC(dec)
	if err != message.ErrTypeMismatch {
		t.Errorf("want ErrTypeMismatch but got %s", err.Error())
	}
}

func TestArray_MarshalELRPC(t *testing.T) {
	v := message.Array[*message.Uint8]{
		Items: []*message.Uint8{
			{Value: 0xaa},
			{Value: 0xbb},
		},
	}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x0a,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // length
		0x01, 0xaa, // [0]: uint8 0xaa
		0x01, 0xbb, // [1]: uint8 0xbb
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestArray_ZeroMessage(t *testing.T) {
	var v message.Array[message.Void]
	got := v.ZeroMessage()
	if _, ok := got.(*message.Array[message.Void]); !ok {
		t.Errorf("want Array but got %T", got)
	}
}

func TestAny_UnmarshalELRPC(t *testing.T) {
	buf := []byte{
		0x0c,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // length
		0x02, 0xab, 0xcd, // uint16 0xabcd
	}
	dec := message.NewDecoder(buf)
	var v message.Any
	err := v.UnmarshalELRPC(dec)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0xab, 0xcd}
	got := v.Raw
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestAny_MarshalELRPC(t *testing.T) {
	v := message.Any{Raw: []byte{0x02, 0xab, 0xcd}}
	enc := message.NewEncoder()
	err := v.MarshalELRPC(enc)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x0c,                                           // type tag
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // length
		0x02, 0xab, 0xcd, // uint16 0xabcd
	}
	got := enc.Buffer()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}

func TestAny_ZeroMessage(t *testing.T) {
	var v message.Any
	got := v.ZeroMessage()
	if _, ok := got.(*message.Any); !ok {
		t.Errorf("want Any but got %T", got)
	}
}
