package x

import (
	"github.com/genkami/elsi/elrpc/helpers"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

type Greeter interface {
	Greet(name *message.Bytes) (*message.Bytes, error)
}

type greeterClient struct {
	greetImpl *helpers.MethodCaller1[*message.Bytes, *message.Bytes]
}

var _ Greeter = &greeterClient{}

func ExportGreeter(rt types.Runtime) Greeter {
	return &greeterClient{
		greetImpl: helpers.NewMethodCaller1[*message.Bytes, *message.Bytes](rt, ModuleID, MethodID_Greeter_Greet),
	}
}

func (g *greeterClient) Greet(name *message.Bytes) (*message.Bytes, error) {
	return g.greetImpl.Call(name)
}
