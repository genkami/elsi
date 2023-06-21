package main

import (
	"fmt"
	"os"

	"github.com/genkami/elsi/elrpc"
	"github.com/genkami/elsi/elrpc/api/x"
)

var theWorld *elrpc.World

func init() {
	w := elrpc.NewWorld()
	handlers := map[string]elrpc.AnyHandler{
		"elsi.x.ping":       &x.PingHandler,
		"elsi.x.add":        &x.AddHandler,
		"elsi.x.div":        &x.DivHandler,
		"elsi.x.write_file": &x.WriteFileHandler,
	}
	for name, h := range handlers {
		w.Register(name, h)
	}
	theWorld = w
}

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
	instance.Use(theWorld)
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
