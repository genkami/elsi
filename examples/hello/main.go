package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	wbuf := []byte("\x04Ping")
	_, err := os.Stdout.Write(wbuf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "req error: %s\n", err.Error())
		os.Exit(1)
	}

	rbuf := make([]byte, 256)
	_, err = io.ReadFull(os.Stdin, rbuf[:1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error (size): %s\n", err.Error())
		os.Exit(1)
	}

	size := int(rbuf[0])
	_, err = io.ReadFull(os.Stdin, rbuf[:size])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error (body): %s\n", err.Error())
		os.Exit(1)
	}

	if !bytes.Equal(rbuf[:size], []byte("Pong")) {
		fmt.Fprintf(os.Stderr, "response mismatch: %v\n", rbuf[:size])
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "hello: OK\n")
}
