package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/genkami/elsi/elrpc/message"
)

var methodsMap = map[string]func(*message.Decoder) ([]byte, error){
	"elsi.x.ping":       ping,
	"elsi.x.add":        add,
	"elsi.x.div":        div,
	"elsi.x.write_file": writeFile,
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

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		strm := &pipeStream{Writer: stdin, Reader: stdout}
		err := serverWorker(strm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worker error: %s\n", err.Error())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "esotime: waiting worker goroutine to finish...\n")
	wg.Wait()

	fmt.Fprintf(os.Stderr, "esotime: OK\n")
}

type Stream interface {
	io.ReadWriter
}

type pipeStream struct {
	io.Reader
	io.Writer
}

func serverWorker(strm Stream) error {
	var err error
	for {
		rlenBuf := make([]byte, message.LengthSize)
		_, err = io.ReadFull(strm, rlenBuf)
		if err != nil {
			return err
		}
		length, err := message.DecodeLength(rlenBuf)
		if err != nil {
			return err
		}

		req := make([]byte, length)
		_, err = io.ReadFull(strm, req)
		if err != nil {
			return err
		}
		dec := message.NewDecoder(req)

		resp, err := dispatchRequest(dec)
		if err != nil {
			return err
		}

		wlenBuf, err := message.AppendLength(nil, len(resp))
		if err != nil {
			return err
		}
		_, err = strm.Write(wlenBuf)
		if err != nil {
			return err
		}

		_, err = strm.Write(resp)
		if err != nil {
			return err
		}
	}
}

func dispatchRequest(dec *message.Decoder) ([]byte, error) {
	methodName, err := dec.DecodeBytes()
	if err != nil {
		return nil, err
	}
	method, ok := methodsMap[string(methodName)]
	if !ok {
		return nil, fmt.Errorf("no such method: %s", string(methodName))
	}
	return method(dec)
}

// ping(nonce: int64) -> int64
func ping(dec *message.Decoder) ([]byte, error) {
	enc := message.NewEncoder()
	nonce, err := dec.DecodeInt64()
	if err != nil {
		return nil, err
	}
	err = enc.EncodeInt64(nonce)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}

// add(x: int64, y: int64) -> int64
func add(dec *message.Decoder) ([]byte, error) {
	enc := message.NewEncoder()
	x, err := dec.DecodeInt64()
	if err != nil {
		return nil, err
	}
	y, err := dec.DecodeInt64()
	if err != nil {
		return nil, err
	}
	err = enc.EncodeInt64(x + y)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}

// div(x: int32, y: int32) -> result: int64 | error: uint64
func div(dec *message.Decoder) ([]byte, error) {
	enc := message.NewEncoder()
	x, err := dec.DecodeInt64()
	if err != nil {
		return nil, err
	}
	y, err := dec.DecodeInt64()
	if err != nil {
		return nil, err
	}
	if y == 0 {
		err = enc.EncodeVariant(1)
		if err != nil {
			return nil, err
		}
		// ZeroDivisionError
		err = enc.EncodeUint64(0xababcdcd)
		if err != nil {
			return nil, err
		}
		return enc.Buffer(), nil
	}
	err = enc.EncodeVariant(0)
	if err != nil {
		return nil, err
	}
	err = enc.EncodeInt64(x / y)
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}

// write_file(handle: uint64, buf: bytes) -> nwritten: uint64 | error: uint64
func writeFile(dec *message.Decoder) ([]byte, error) {
	enc := message.NewEncoder()
	_, err := dec.DecodeUint64()
	if err != nil {
		return nil, err
	}
	buf, err := dec.DecodeBytes()
	if err != nil {
		return nil, err
	}

	nwritten, err := os.Stdout.Write(buf)
	if err != nil {
		return nil, err
	}

	err = enc.EncodeVariant(0)
	if err != nil {
		return nil, err
	}
	err = enc.EncodeUint64(uint64(nwritten))
	if err != nil {
		return nil, err
	}
	return enc.Buffer(), nil
}
