package apibuilder_test

import (
	"testing"

	"github.com/genkami/elsi/elrpc/apibuilder"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
	"github.com/google/go-cmp/cmp"
)

func TestHostHandler0(t *testing.T) {
	enc := message.NewEncoder()

	type Handler = apibuilder.HostHandler0[*message.String]
	handler := Handler(func() (*message.String, error) {
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func TestHostHandler1(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(123)
	if err != nil {
		t.Fatal(err)
	}

	type Handler = apibuilder.HostHandler1[*message.Uint32, *message.String]
	handler := Handler(func(x1 *message.Uint32) (*message.String, error) {
		if x1.Value != 123 {
			t.Errorf("want 123 but got %d", x1.Value)
		}
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func TestHostHandler2(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(123)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeInt16(-45)
	if err != nil {
		t.Fatal(err)
	}

	type Handler = apibuilder.HostHandler2[*message.Uint32, *message.Int16, *message.String]
	handler := Handler(func(x1 *message.Uint32, x2 *message.Int16) (*message.String, error) {
		if x1.Value != 123 {
			t.Errorf("want 123 but got %d", x1.Value)
		}
		if x2.Value != -45 {
			t.Errorf("want -45 but got %d", x2.Value)
		}
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func TestHostHandler3(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(123)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeInt16(-45)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeString("67")
	if err != nil {
		t.Fatal(err)
	}

	type Handler = apibuilder.HostHandler3[*message.Uint32, *message.Int16, *message.String, *message.String]
	handler := Handler(func(x1 *message.Uint32, x2 *message.Int16, x3 *message.String) (*message.String, error) {
		if x1.Value != 123 {
			t.Errorf("want 123 but got %d", x1.Value)
		}
		if x2.Value != -45 {
			t.Errorf("want -45 but got %d", x2.Value)
		}
		if x3.Value != "67" {
			t.Errorf("want \"67\" but got %q", x3.Value)
		}
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func TestHostHandler4(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(123)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeInt16(-45)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeString("67")
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeUint8(8)
	if err != nil {
		t.Fatal(err)
	}

	type Handler = apibuilder.HostHandler4[*message.Uint32, *message.Int16, *message.String, *message.Uint8, *message.String]
	handler := Handler(func(x1 *message.Uint32, x2 *message.Int16, x3 *message.String, x4 *message.Uint8) (*message.String, error) {
		if x1.Value != 123 {
			t.Errorf("want 123 but got %d", x1.Value)
		}
		if x2.Value != -45 {
			t.Errorf("want -45 but got %d", x2.Value)
		}
		if x3.Value != "67" {
			t.Errorf("want \"67\" but got %q", x3.Value)
		}
		if x4.Value != 8 {
			t.Errorf("want 8 but got %d", x4.Value)
		}
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func TestHostHandler5(t *testing.T) {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(123)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeInt16(-45)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeString("67")
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeUint8(8)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeInt64(-90)
	if err != nil {
		t.Fatal(err)
	}

	type Handler = apibuilder.HostHandler5[*message.Uint32, *message.Int16, *message.String, *message.Uint8, *message.Int64, *message.String]
	handler := Handler(func(x1 *message.Uint32, x2 *message.Int16, x3 *message.String, x4 *message.Uint8, x5 *message.Int64) (*message.String, error) {
		if x1.Value != 123 {
			t.Errorf("want 123 but got %d", x1.Value)
		}
		if x2.Value != -45 {
			t.Errorf("want -45 but got %d", x2.Value)
		}
		if x3.Value != "67" {
			t.Errorf("want \"67\" but got %q", x3.Value)
		}
		if x4.Value != 8 {
			t.Errorf("want 8 but got %d", x4.Value)
		}
		if x5.Value != -90 {
			t.Errorf("want -90 but got %d", x5.Value)
		}
		return &message.String{Value: "hello"}, nil
	})

	got, err := handler.HandleRequest(message.NewDecoder(enc.Buffer()))
	if err != nil {
		t.Fatal(err)
	}

	want := &message.String{Value: "hello"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

type mockRuntime struct {
	ret  *message.Any
	err  error
	call runtimeCall
}

type runtimeCall struct {
	ModuleID uint32
	MethodID uint32
	Args     *message.Any
}

var _ types.Runtime = (*mockRuntime)(nil)

func (rt *mockRuntime) Use(moduleID uint32, methodID uint32, handler types.HostHandler) {
	// nop
}

func (rt *mockRuntime) Call(moduleID uint32, methodID uint32, args *message.Any) (*message.Any, error) {
	rt.call = runtimeCall{
		ModuleID: moduleID,
		MethodID: methodID,
		Args:     args,
	}
	return rt.ret, rt.err
}

func TestGuestDelegator0(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator0[*message.String](rt, modID, methodID)
	gotRet, err := delegator.Call()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args:     encodeAsAny(t),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func TestGuestDelegator1(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator1[*message.Uint8, *message.String](rt, modID, methodID)
	gotRet, err := delegator.Call(
		&message.Uint8{Value: 123},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args: encodeAsAny(
			t,
			&message.Uint8{Value: 123},
		),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func TestGuestDelegator2(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator2[*message.Uint8, *message.Int16, *message.String](rt, modID, methodID)
	gotRet, err := delegator.Call(
		&message.Uint8{Value: 123},
		&message.Int16{Value: -45},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args: encodeAsAny(
			t,
			&message.Uint8{Value: 123},
			&message.Int16{Value: -45},
		),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func TestGuestDelegator3(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator3[*message.Uint8, *message.Int16, *message.Uint32, *message.String](rt, modID, methodID)
	gotRet, err := delegator.Call(
		&message.Uint8{Value: 123},
		&message.Int16{Value: -45},
		&message.Uint32{Value: 67},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args: encodeAsAny(
			t,
			&message.Uint8{Value: 123},
			&message.Int16{Value: -45},
			&message.Uint32{Value: 67},
		),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func TestGuestDelegator4(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator4[*message.Uint8, *message.Int16, *message.Uint32, *message.Int64, *message.String](rt, modID, methodID)
	gotRet, err := delegator.Call(
		&message.Uint8{Value: 123},
		&message.Int16{Value: -45},
		&message.Uint32{Value: 67},
		&message.Int64{Value: -8},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args: encodeAsAny(
			t,
			&message.Uint8{Value: 123},
			&message.Int16{Value: -45},
			&message.Uint32{Value: 67},
			&message.Int64{Value: -8},
		),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func TestGuestDelegator5(t *testing.T) {
	modID := uint32(0x0000_1234)
	methodID := uint32(0x0000_5678)

	wantRet := &message.String{Value: "Ok"}
	rt := &mockRuntime{
		ret: encodeAsAny(t, wantRet),
	}

	delegator := apibuilder.NewGuestDelegator5[*message.Uint8, *message.Int16, *message.Uint32, *message.Int64, *message.String, *message.String](rt, modID, methodID)
	gotRet, err := delegator.Call(
		&message.Uint8{Value: 123},
		&message.Int16{Value: -45},
		&message.Uint32{Value: 67},
		&message.Int64{Value: -8},
		&message.String{Value: "abc"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("mismatch (-want +got)\n%s", diff)
	}

	wantCall := runtimeCall{
		ModuleID: modID,
		MethodID: methodID,
		Args: encodeAsAny(
			t,
			&message.Uint8{Value: 123},
			&message.Int16{Value: -45},
			&message.Uint32{Value: 67},
			&message.Int64{Value: -8},
			&message.String{Value: "abc"},
		),
	}
	if diff := cmp.Diff(wantCall, rt.call); diff != "" {
		t.Errorf("mismtach (-want +got)\n%s", diff)
	}
}

func encodeAsAny(t *testing.T, msgs ...message.Message) *message.Any {
	enc := message.NewEncoder()
	for _, msg := range msgs {
		err := msg.MarshalELRPC(enc)
		if err != nil {
			t.Fatal(err)
		}
	}
	return &message.Any{Raw: enc.Buffer()}
}
