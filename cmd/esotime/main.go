package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/genkami/elsi/elrpc/runtime"
	"github.com/genkami/elsi/elsi/api/exp"
	"github.com/genkami/elsi/elsi/impl/expimpl"
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

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	guest := runtime.NewProcessGuest(args[2], args[3:]...)
	rt := runtime.NewRuntime(logger, guest)

	hs := expimpl.NewHandleSet()
	stdio := expimpl.NewStdio(hs, map[uint8]expimpl.StdHandleCtor{
		exp.HandleTypeStdin: func() (any, error) {
			return os.Stdin, nil
		},
		exp.HandleTypeStdout: func() (any, error) {
			return os.Stdout, nil
		},
		exp.HandleTypeStderr: func() (any, error) {
			return os.Stderr, nil
		},
	})
	exp.ImportStdio(rt, stdio)

	stream := expimpl.NewStream(hs)
	exp.ImportStream(rt, stream)

	file := expimpl.NewFile(hs)
	exp.ImportFile(rt, file)

	http := expimpl.NewHTTP(logger, hs, map[string]expimpl.HttpListenerConfig{
		"default": {
			AddrAndPort: ":8080",
		},
	})
	exp.ImportHTTP(rt, http)

	// TODO: use UseWorld

	err := rt.Start()
	if err != nil {
		panic(err)
	}

	err = rt.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}
