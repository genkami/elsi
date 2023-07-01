package main

import (
	"fmt"
	"io"
	"os"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
)

func main() {
	actions := []func() error{
		doPing,
		doAdd,
		doDiv(14, 7),
		doDiv(15, 0),
		doWriteFile,
		doTestExport,
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
	err := enc.EncodeUint64(0x0000_BEEF_0000_0000) // elsi.x.TODO/Ping
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

	vt, err := dec.DecodeVariantTag()
	switch vt {
	case 0:
		gotNonce, err := dec.DecodeInt64()
		if err != nil {
			return err
		}

		if nonce != gotNonce {
			return fmt.Errorf("nonce mismatch: want = %d, got = %d", nonce, gotNonce)
		}
	case 1:
		code, err := dec.DecodeUint64()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "code: error (code = %X)\n", code)
	default:
		return fmt.Errorf("unknown variant: %d", vt)
	}

	fmt.Fprintf(os.Stderr, "ping: OK\n")
	return nil
}

func doAdd() error {
	var err error
	var x int64 = 333
	var y int64 = 222

	enc := message.NewEncoder()
	err = enc.EncodeUint64(0x0000_BEEF_0000_0001) // elsi.x.TODO/Add
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

	vt, err := dec.DecodeVariantTag()
	if err != nil {
		return err
	}
	switch vt {
	case 0:
		sum, err := dec.DecodeInt64()
		if err != nil {
			return err
		}

		if sum != x+y {
			return fmt.Errorf("%d + %d should be %d but got %d", x, y, x+y, sum)
		}
	case 1:
		code, err := dec.DecodeUint64()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "code: error (code = %X)\n", code)
	default:
		return fmt.Errorf("unknown variant: %d", vt)
	}
	fmt.Fprintf(os.Stderr, "add: OK\n")
	return nil
}

func doDiv(x, y int64) func() error {
	return func() error {
		var err error

		enc := message.NewEncoder()
		err = enc.EncodeUint64(0x0000_BEEF_0000_0002) // elsi.x.TODO/div
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

		vtag, err := dec.DecodeVariantTag()
		if err != nil {
			return err
		}

		switch vtag {
		case 0:
			sum, err := dec.DecodeInt64()
			if err != nil {
				return err
			}

			if sum != x/y {
				return fmt.Errorf("%d / %d should be %d but got %d", x, y, x/y, sum)
			}
		case 1:
			code, err := dec.DecodeUint64()
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "div: error (code = %X)\n", code)
		default:
			return fmt.Errorf("unknown variant: %d", vtag)
		}

		fmt.Fprintf(os.Stderr, "div: OK\n")
		return nil
	}
}

func doWriteFile() error {
	var err error
	enc := message.NewEncoder()
	err = enc.EncodeUint64(0x0000_BEEF_0000_0003) // elsi.x.TODO/WriteFile
	if err != nil {
		return err
	}
	err = enc.EncodeUint64(1)
	if err != nil {
		return err
	}
	err = enc.EncodeBytes([]byte("Hello from ELSI!\n"))
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

	vtag, err := dec.DecodeVariantTag()
	if err != nil {
		return err
	}

	switch vtag {
	case 0:
		nwritten, err := dec.DecodeUint64()
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%d bytes written\n", nwritten)
	case 1:
		code, err := dec.DecodeUint64()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "writeFile: error (code = %X)\n", code)
	default:
		return fmt.Errorf("unknown variant: %d", vtag)
	}

	fmt.Fprintf(os.Stderr, "writeFile: OK\n")
	return nil
}

func doTestExport() error {
	var err error
	enc := message.NewEncoder()
	err = enc.EncodeUint64(0x0000_BEEF_0000_0004) // elsi.x.TODO/TestExport
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

	vtag, err := dec.DecodeVariantTag()
	if err != nil {
		return err
	}

	switch vtag {
	case 0:
		// OK (nop)
	case 1:
		code, err := dec.DecodeUint64()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "testExport: error (code = %X)\n", code)
	default:
		return fmt.Errorf("unknown variant: %d", vtag)
	}

	var callID uint64
	var name string

pollLoop:
	for {
		enc = message.NewEncoder()
		err = enc.EncodeUint64(0x0000_0000_0000_0000) // builtin.PollMethodCall
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

		vtag, err := dec.DecodeVariantTag()
		if err != nil {
			return err
		}
		switch vtag {
		case 0:
			callID, err = dec.DecodeUint64()
			if err != nil {
				return err
			}
			mID, err := dec.DecodeUint64()
			if err != nil {
				return err
			}
			if mID != 0x0000_BEEF_0000_0010 { // elsi.x.Greeter/Greet
				return fmt.Errorf("unknown method ID: %X", mID)
			}
			arg, err := dec.DecodeAny()
			if err != nil {
				return err
			}
			dec = message.NewDecoder(arg.Raw)
			name, err = dec.DecodeString()
			if err != nil {
				return err
			}
			break pollLoop
		case 1:
			code, err := dec.DecodeUint64()
			if err != nil {
				return err
			}
			if code != builtin.CodeNotFound {
				fmt.Fprintf(os.Stderr, "testExport: error (code = %X)\n", code)
			}
		default:
			return fmt.Errorf("unknown variant: %d", vtag)
		}
	}

	enc = message.NewEncoder()
	err = enc.EncodeUint64(0x0000_0000_0000_0001) // builtin.SendResult
	if err != nil {
		return err
	}
	err = enc.EncodeUint64(callID)
	if err != nil {
		return err
	}
	err = enc.EncodeVariantTag(0)
	if err != nil {
		return err
	}
	anyEnc := message.NewEncoder()
	err = anyEnc.EncodeString("Hello, " + name + "!")
	if err != nil {
		return err
	}
	err = enc.EncodeAny(&message.Any{Raw: anyEnc.Buffer()})
	if err != nil {
		return err
	}

	err = sendReq(enc.Buffer())
	if err != nil {
		return err
	}

	dec, err = receiveResp()
	if err != nil {
		return err
	}

	vtag, err = dec.DecodeVariantTag()
	if err != nil {
		return err
	}
	switch vtag {
	case 0:
		// OK (nop)
	case 1:
		code, err := dec.DecodeUint64()
		if err != nil {
			return err
		}
		msg, err := dec.DecodeString()
		if err != nil {
			return err
		}
		return fmt.Errorf("testExport: SendResult: error (code = %X): %s", code, msg)
	default:
		return fmt.Errorf("testExport: SendResult: unknown variant: %d", vtag)
	}

	fmt.Fprintf(os.Stderr, "testExport: OK\n")
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
