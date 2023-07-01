package builtinimpl_test

import (
	"os"
	"testing"
	"time"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/message/msgtest"
	"github.com/genkami/elsi/elrpc/runtime/internal/builtinimpl"
	"golang.org/x/exp/slog"
)

const (
	ModuleID = 0x0000_FFFF // only used by test package

	MethodID_Nop_Nop = 0x0000_0000
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func TestExporterImpl_PollMethodCall_notfound(t *testing.T) {
	e := builtinimpl.NewExporter(logger)
	_, err := e.PollMethodCall()
	if err == nil {
		t.Errorf("want error but got nil")
	}
	msgtest.AssertError(t, err, builtin.ModuleID, builtin.CodeNotFound)
}

func TestExporterImpl_PollMethodCall_found(t *testing.T) {
	e := builtinimpl.NewExporter(logger)

	enc := message.NewEncoder()
	err := enc.EncodeUint64(0xdeadbeef)
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
	if got.ModuleID != ModuleID {
		t.Errorf("want ModuleID 0x%X but got 0x%X", ModuleID, got.ModuleID)
	}
	if got.MethodID != MethodID_Nop_Nop {
		t.Errorf("want MethodID 0x%X but got 0x%X", MethodID_Nop_Nop, got.MethodID)
	}

	dec := message.NewDecoder(got.Args.Raw)
	n, err := dec.DecodeUint64()
	if err != nil {
		t.Fatal(err)
	}
	if n != 0xdeadbeef {
		t.Errorf("want 0xdeadbeef but got 0x%X", n)
	}

	select {
	case r := <-resultCh:
		t.Fatalf("result should not be sent yet, but got %#v", r)
	default:
	}

	enc = message.NewEncoder()
	err = enc.EncodeString("OK")
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
			t.Errorf("expected OK but got error: %s", r.RetVal.Err.Error())
		}
		dec := message.NewDecoder(r.RetVal.Ok.Raw)
		s, err := dec.DecodeString()
		if err != nil {
			t.Fatal(err)
		}
		if s != "OK" {
			t.Errorf("want OK but got %s", s)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
