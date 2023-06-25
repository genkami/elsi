package main

import (
	"fmt"
	"os"

	"github.com/genkami/elsi/elrpc"
	"github.com/genkami/elsi/elrpc/api/x"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: esotime run CMD...\n")
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}
	if args[1] != "run" {
		usage()
	}

	mod := elrpc.NewProcessModule(args[2], args[3:]...)
	instance := elrpc.NewInstance(mod)
	instance.Use(x.NewWorld(&todoImpl{}))
	err := instance.Start()
	if err != nil {
		panic(err)
	}

	err = instance.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}

type todoImpl struct{}

var _ x.TODO = &todoImpl{}

func (*todoImpl) Ping(req *x.PingRequest) (*x.PingResponse, error) {
	return &x.PingResponse{
		Nonce: req.Nonce,
	}, nil
}

func (*todoImpl) Add(req *x.AddRequest) (*x.AddResponse, error) {
	return &x.AddResponse{
		Sum: req.X + req.Y,
	}, nil
}

func (*todoImpl) Div(req *x.DivRequest) (*x.DivResponse, error) {
	type Resp = elrpc.Result[*x.DivResponse, *elrpc.Error]
	if req.Y == 0 {
		return nil, &elrpc.Error{
			Code:    0xdeadbeef,
			Message: "division by zero",
		}
	}
	return &x.DivResponse{
		Result: req.X / req.Y,
	}, nil
}

func (*todoImpl) WriteFile(req *x.WriteFileRequest) (*x.WriteFileResponse, error) {
	type Resp = elrpc.Result[*x.WriteFileResponse, *elrpc.Error]
	length, err := os.Stdout.Write(req.Buf)
	if err != nil {
		return nil, err
	}

	return &x.WriteFileResponse{
		Length: uint64(length),
	}, nil
}
