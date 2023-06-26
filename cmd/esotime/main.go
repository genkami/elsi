package main

import (
	"fmt"
	"os"
	"time"

	"github.com/genkami/elsi/elrpc/api/x"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/runtime"
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

	mod := runtime.NewProcessModule(args[2], args[3:]...)
	instance := runtime.NewInstance(mod)
	todo := &todoImpl{}
	exports := x.UseWorld(instance, &x.Imports{
		TODO: todo,
	})
	todo.greeter = exports.Greeter
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

type todoImpl struct {
	greeter x.Greeter
}

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
	if req.Y == 0 {
		return nil, &message.Error{
			Code:    0xdeadbeef,
			Message: "division by zero",
		}
	}
	return &x.DivResponse{
		Result: req.X / req.Y,
	}, nil
}

func (*todoImpl) WriteFile(req *x.WriteFileRequest) (*x.WriteFileResponse, error) {
	length, err := os.Stdout.Write(req.Buf)
	if err != nil {
		return nil, err
	}

	return &x.WriteFileResponse{
		Length: uint64(length),
	}, nil
}

func (t *todoImpl) TestExport() (*message.Void, error) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		res, err := t.greeter.Greet(&message.Bytes{Value: []byte("Taro")})
		if err != nil {
			fmt.Fprintf(os.Stderr, "esotime: TestExport: %s\n", err.Error())
			return
		}
		fmt.Fprintf(os.Stderr, "esotime: GreetResponse = %s\n", string(res.Value))
	}()
	return &message.Void{}, nil
}
