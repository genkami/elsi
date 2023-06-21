package x

import "github.com/genkami/elsi/elrpc"

type Greeter interface {
	Greet(name *elrpc.Bytes) *elrpc.Bytes
}

type GreeterClient struct {
	instance *elrpc.Instance
}

var _ Greeter = &GreeterClient{}

// TODO: handler functions should return error
// -> Or we can let every methods have Result<_, Error> as a return value.
// TODO: elrpc.Handler should have HandleImport and HandleExport
func (g *GreeterClient) Greet(name *elrpc.Bytes) *elrpc.Bytes {
	enc := elrpc.NewEncoder()
	err := name.MarshalELRPC(enc)
	if err != nil {
		panic("TODO")
	}
	rawResp := g.instance.Call([]byte("elsi.x.greeter/greet"), &elrpc.Any{Raw: enc.Buffer()})
	dec := elrpc.NewDecoder(rawResp.Raw)
	resp := elrpc.NewMessage[*elrpc.Bytes]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		panic("TODO")
	}
	return resp.(*elrpc.Bytes)
}
