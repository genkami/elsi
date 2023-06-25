package x

import "github.com/genkami/elsi/elrpc"

type Greeter interface {
	Greet(name *elrpc.Bytes) (*elrpc.Bytes, error)
}

type greeterClient struct {
	greetImpl *elrpc.MethodCaller1[*elrpc.Bytes, *elrpc.Bytes]
}

var _ Greeter = &greeterClient{}

func ExportGreeter(instance *elrpc.Instance) Greeter {
	return &greeterClient{
		greetImpl: elrpc.NewMethodCaller1[*elrpc.Bytes, *elrpc.Bytes](instance, ModuleID, MethodID_Greeter_Greet),
	}
}

func (g *greeterClient) Greet(name *elrpc.Bytes) (*elrpc.Bytes, error) {
	return g.greetImpl.Call(name)
}
