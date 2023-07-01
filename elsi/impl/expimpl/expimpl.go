package expimpl

import (
	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
)

var (
	errNoSuchHandle = &message.Error{
		ModuleID: builtin.ModuleID,
		Code:     builtin.CodeNotFound,
		Message:  "no such handle",
	}
	errInvalidHandleType = &message.Error{
		ModuleID: builtin.ModuleID,
		Code:     builtin.CodeNotFound,
		Message:  "invalid handle type",
	}
	errUnsupported = &message.Error{
		ModuleID: exp.ModuleID,
		Code:     exp.CodeUnsupported,
		Message:  "unsupported operation",
	}
	errNoRequest = &message.Error{
		ModuleID: builtin.ModuleID,
		Code:     builtin.CodeNotFound,
		Message:  "no request",
	}
)
