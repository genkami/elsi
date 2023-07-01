package main

import (
	"io"
	"os"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
)

func main() {
	err := sendReq(
		exp.ModuleID, exp.MethodID_Stdio_OpenStdHandle,
		&message.Uint8{Value: exp.HandleTypeStdout},
	)
	if err != nil {
		panic(err)
	}
	dec, err := receiveResp()
	if err != nil {
		panic(err)
	}
	stdout := &message.Result[*exp.Handle, *message.Error]{}
	err = stdout.UnmarshalELRPC(dec)
	if err != nil {
		panic(err)
	}
	if !stdout.IsOk {
		panic(stdout.Err)
	}

	err = sendReq(
		exp.ModuleID, exp.MethodID_Stream_Write,
		stdout.Ok, &message.Bytes{Value: []byte("Hello, world!\n")},
	)
	if err != nil {
		panic(err)
	}
	dec, err = receiveResp()
	if err != nil {
		panic(err)
	}
	writeResult := &message.Result[*message.Uint64, *message.Error]{}
	err = writeResult.UnmarshalELRPC(dec)
	if err != nil {
		panic(err)
	}
	if !writeResult.IsOk {
		panic(writeResult.Err)
	}
}

func sendReq(modID, methodID uint32, args ...message.Message) error {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(modID)
	if err != nil {
		return err
	}
	err = enc.EncodeUint32(methodID)
	if err != nil {
		return err
	}
	for _, arg := range args {
		err = arg.MarshalELRPC(enc)
		if err != nil {
			return err
		}
	}

	buf := enc.Buffer()
	lenBuf, err := message.AppendLength(nil, len(buf))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(lenBuf)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(buf)
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
	msgLen, err := message.DecodeLength(lenBuf)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, msgLen)
	_, err = io.ReadFull(os.Stdin, buf)
	if err != nil {
		return nil, err
	}
	return message.NewDecoder(buf), nil
}
