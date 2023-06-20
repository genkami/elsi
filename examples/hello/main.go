package main

import (
	"fmt"
	"io"
	"os"

	"github.com/genkami/elsi/elrpc/message"
)

func main() {
	actions := []func() error{
		doPing,
		doAdd,
	}

	for _, action := range actions {
		err := action()
		if err != nil {
			fmt.Fprintf(os.Stderr, "action error: %s\n", err.Error())
		}
	}

	fmt.Fprintf(os.Stderr, "hello: OK\n")
}

func doPing() error {
	enc := message.NewEncoder()
	err := enc.EncodeBytes([]byte("elsi.x.ping"))
	if err != nil {
		return err
	}
	var nonce int64 = 12345
	err = enc.EncodeInt64(nonce)
	if err != nil {
		return err
	}

	err = sendReq(enc.Buffer())
	if err != nil {
		return err
	}
	dec, err := receiveResp()
	if err != nil {
		return err
	}

	gotNonce, err := dec.DecodeInt64()
	if err != nil {
		return err
	}

	if nonce != gotNonce {
		return fmt.Errorf("nonce mismatch: want = %d, got = %d", nonce, gotNonce)
	}
	fmt.Fprintf(os.Stderr, "ping: OK\n")
	return nil
}

func doAdd() error {
	var err error
	var x int64 = 333
	var y int64 = 222

	enc := message.NewEncoder()
	err = enc.EncodeBytes([]byte("elsi.x.add"))
	if err != nil {
		return err
	}
	err = enc.EncodeInt64(x)
	if err != nil {
		return err
	}
	err = enc.EncodeInt64(y)
	if err != nil {
		return err
	}

	err = sendReq(enc.Buffer())
	if err != nil {
		return err
	}
	dec, err := receiveResp()
	if err != nil {
		return err
	}

	sum, err := dec.DecodeInt64()
	if err != nil {
		return err
	}

	if sum != x+y {
		return fmt.Errorf("%d + %d should be %d but got %d", x, y, x+y, sum)
	}
	fmt.Fprintf(os.Stderr, "add: OK\n")
	return nil
}

func sendReq(req []byte) error {
	lenBuf, err := message.AppendLength(nil, len(req))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(lenBuf)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(req)
	if err != nil {
		return err
	}
	return nil
}

func receiveResp() (*message.Decoder, error) {
	lenBuf := make([]byte, message.LengthSize)
	_, err := io.ReadFull(os.Stdin, lenBuf)
	if err != nil {
		return nil, err
	}
	length, err := message.DecodeLength(lenBuf)
	if err != nil {
		return nil, err
	}
	resp := make([]byte, length)
	_, err = io.ReadFull(os.Stdin, resp)
	if err != nil {
		return nil, err
	}
	return message.NewDecoder(resp), nil
}
