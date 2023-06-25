package x

import "github.com/genkami/elsi/elrpc"

type Greeter interface {
	Greet(name *elrpc.Bytes) (*elrpc.Bytes, error)
}

type greeterClient struct {
	instance *elrpc.Instance
}

var _ Greeter = &greeterClient{}

func ExportGreeter(instance *elrpc.Instance) Greeter {
	return &greeterClient{instance: instance}
}

// TODO: handler functions should return error
// -> Or we can let every methods have Result<_, Error> as a return value.
// TODO: elrpc.Handler should have HandleImport and HandleExport
func (g *greeterClient) Greet(name *elrpc.Bytes) (*elrpc.Bytes, error) {
	enc := elrpc.NewEncoder()
	err := name.MarshalELRPC(enc)
	if err != nil {
		return nil, err
	}
	rawResp, err := g.instance.Call([]byte("elsi.x.greeter/greet"), &elrpc.Any{Raw: enc.Buffer()})
	if err != nil {
		return nil, err
	}
	dec := elrpc.NewDecoder(rawResp.Raw)
	resp := &elrpc.Bytes{}
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
