package msgtest

import (
	"errors"
	"testing"

	"github.com/genkami/elsi/elrpc/message"
)

func AssertError(t *testing.T, err error, modID, code uint32) {
	t.Helper()
	var msgErr *message.Error
	if !errors.As(err, &msgErr) {
		t.Errorf("want *message.Error but got %T (%s)", err, err.Error())
		return
	}
	if msgErr.ModuleID != modID || msgErr.Code != code {
		t.Errorf("want error (mod = %X, code = %X) but got %s", modID, code, err.Error())
	}
}
