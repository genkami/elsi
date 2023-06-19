package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	cmd := exec.Command(args[2], args[3:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 256)
	_, err = io.ReadFull(stdout, buf[:1])
	if err != nil {
		panic(err)
	}
	size := int(buf[0])
	_, err = io.ReadFull(stdout, buf[:size])
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf[:size], []byte("Ping")) {
		panic(fmt.Sprintf("Protocol error: %v", buf[:size]))
	}

	wbuf := []byte("\x04Pong")
	_, err = stdin.Write(wbuf)
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}
