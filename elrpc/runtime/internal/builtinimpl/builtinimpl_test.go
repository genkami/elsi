package builtinimpl_test

import (
	"os"
	"testing"
	"time"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/elrpctest"
	"github.com/genkami/elsi/elrpc/runtime/internal/builtinimpl"
	"golang.org/x/exp/slog"
)

const (
	ModuleID = 0x0000_ffff // only used by test package

	MethodID_Nop_Nop = 0x0000_abcd

	CodeFoo = 0x0000_1234
)

var (
	logger  = slog.New(slog.NewTextHandler(os.Stderr, nil))
	timeout = 1 * time.Second
)

func TestExporterImpl_PollMethodCall_notfound(t *testing.T) {
	e := builtinimpl.NewExporter(logger)
	_, err := e.PollMethodCall()
	elrpctest.AssertError(t, err, builtin.ModuleID, builtin.CodeNotFound)
}

func TestExporterImpl_PollMethodCall_found(t *testing.T) {
	e := builtinimpl.NewExporter(logger)

	enc := message.NewEncoder()
	err := enc.EncodeString("Ping")
	if err != nil {
		t.Fatal(err)
	}

	call := &builtin.MethodCall{
		ModuleID: ModuleID,
		MethodID: MethodID_Nop_Nop,
		Args:     &message.Any{Raw: enc.Buffer()},
	}
	resultCh := e.CallAsync(call)

	got, err := e.PollMethodCall()
	if err != nil {
		t.Fatal(err)
	}
	if got.ModuleID != ModuleID || got.MethodID != MethodID_Nop_Nop {
		t.Errorf("want (mod = %X, method = %X) but got (mod = %X, method = %X)",
			ModuleID, MethodID_Nop_Nop, got.ModuleID, got.MethodID)
	}

	dec := message.NewDecoder(got.Args.Raw)
	resp, err := dec.DecodeString()
	if err != nil {
		t.Fatal(err)
	}
	if resp != "Ping" {
		t.Errorf("want Ping but got %s", resp)
	}

	select {
	case r := <-resultCh:
		t.Fatalf("result should not be sent yet, but got %#v", r)
	default:
	}
}

func TestExporterImpl_SendResult_ok(t *testing.T) {
	e := builtinimpl.NewExporter(logger)

	enc := message.NewEncoder()
	err := enc.EncodeString("Ping")
	if err != nil {
		t.Fatal(err)
	}

	call := &builtin.MethodCall{
		ModuleID: ModuleID,
		MethodID: MethodID_Nop_Nop,
		Args:     &message.Any{Raw: enc.Buffer()},
	}
	resultCh := e.CallAsync(call)

	got, err := e.PollMethodCall()
	if err != nil {
		t.Fatal(err)
	}
	if got.ModuleID != ModuleID || got.MethodID != MethodID_Nop_Nop {
		t.Errorf("want (mod = %X, method = %X) but got (mod = %X, method = %X)",
			ModuleID, MethodID_Nop_Nop, got.ModuleID, got.MethodID)
	}

	dec := message.NewDecoder(got.Args.Raw)
	resp, err := dec.DecodeString()
	if err != nil {
		t.Fatal(err)
	}
	if resp != "Ping" {
		t.Errorf("want Ping but got %s", resp)
	}

	select {
	case r := <-resultCh:
		t.Fatalf("result should not be sent yet, but got %#v", r)
	default:
	}

	enc = message.NewEncoder()
	err = enc.EncodeString("Pong")
	if err != nil {
		t.Fatal(err)
	}

	_, err = e.SendResult(&builtin.MethodResult{
		CallID: call.CallID,
		RetVal: &message.Result[*message.Any, *message.Error]{
			IsOk: true,
			Ok:   &message.Any{Raw: enc.Buffer()},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-resultCh:
		if !r.RetVal.IsOk {
			t.Fatal(r.RetVal.Err)
		}
		dec := message.NewDecoder(r.RetVal.Ok.Raw)
		s, err := dec.DecodeString()
		if err != nil {
			t.Fatal(err)
		}
		if s != "Pong" {
			t.Errorf("want Pong but got %s", s)
		}
	case <-time.After(timeout):
		t.Fatal("timeout")
	}
}

func TestExporterImpl_SendResult_err(t *testing.T) {
	e := builtinimpl.NewExporter(logger)

	enc := message.NewEncoder()
	err := enc.EncodeString("Ping")
	if err != nil {
		t.Fatal(err)
	}

	call := &builtin.MethodCall{
		ModuleID: ModuleID,
		MethodID: MethodID_Nop_Nop,
		Args:     &message.Any{Raw: enc.Buffer()},
	}
	resultCh := e.CallAsync(call)

	got, err := e.PollMethodCall()
	if err != nil {
		t.Fatal(err)
	}
	if got.ModuleID != ModuleID || got.MethodID != MethodID_Nop_Nop {
		t.Errorf("want (mod = %X, method = %X) but got (mod = %X, method = %X)",
			ModuleID, MethodID_Nop_Nop, got.ModuleID, got.MethodID)
	}

	dec := message.NewDecoder(got.Args.Raw)
	resp, err := dec.DecodeString()
	if err != nil {
		t.Fatal(err)
	}
	if resp != "Ping" {
		t.Errorf("want Ping but got %s", resp)
	}

	select {
	case r := <-resultCh:
		t.Fatalf("result should not be sent yet, but got %#v", r)
	default:
	}

	enc = message.NewEncoder()
	err = enc.EncodeString("Pong")
	if err != nil {
		t.Fatal(err)
	}

	_, err = e.SendResult(&builtin.MethodResult{
		CallID: call.CallID,
		RetVal: &message.Result[*message.Any, *message.Error]{
			IsOk: false,
			Err: &message.Error{
				ModuleID: ModuleID,
				Code:     CodeFoo,
				Message:  "foo",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-resultCh:
		if r.RetVal.IsOk {
			t.Fatalf("want error but got %#v", r.RetVal.Ok)
		}
		elrpctest.AssertError(t, r.RetVal.Err, ModuleID, CodeFoo)
	case <-time.After(timeout):
		t.Fatal("timeout")
	}
}
