package runtime_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/elrpctest"
	"github.com/genkami/elsi/elrpc/helpers"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/runtime"
	"github.com/genkami/elsi/elrpc/types"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slog"
)

const (
	ModuleID = 0x0000_ffff // only used by test package

	MethodID_HostAPI_Ping         = 0x0000_1234
	MethodID_HostAPI_NoSuchMethod = 0x0000_ffff

	CodeSomeError = 0x0000_abcd
)

type HostAPI interface {
	Ping(*message.String) (*message.String, error)
}

func UseHostAPI(instance types.Instance, h HostAPI) {
	instance.Use(ModuleID, MethodID_HostAPI_Ping, helpers.TypedHandler1[*message.String, *message.String](h.Ping))
}

type hostAPIImpl struct {
	pingImpl func(*message.String) (*message.String, error)
}

func (h *hostAPIImpl) Ping(args *message.String) (*message.String, error) {
	return h.pingImpl(args)
}

func TestInstance_Call(t *testing.T) {
	type Result = message.Result[*message.String, *message.Error]
	cases := []struct {
		name         string
		moduleID     uint32
		methodID     uint32
		req          *message.String
		pingImpl     func(*message.String) (*message.String, error)
		wantOk       *message.String
		wantErrModID uint32
		wantErrCode  uint32
	}{
		{
			name:     "ok",
			moduleID: ModuleID,
			methodID: MethodID_HostAPI_Ping,
			req: &message.String{
				Value: "Ping",
			},
			pingImpl: func(arg *message.String) (*message.String, error) {
				if arg.Value != "Ping" {
					t.Errorf("want Ping but got %s", arg.Value)
					return nil, errors.New("ERROR")
				}
				return &message.String{Value: "Pong"}, nil
			},
			wantOk: &message.String{
				Value: "Pong",
			},
		},
		{
			name:     "elrpc err",
			moduleID: ModuleID,
			methodID: MethodID_HostAPI_Ping,
			req: &message.String{
				Value: "Ping",
			},
			pingImpl: func(arg *message.String) (*message.String, error) {
				return nil, &message.Error{
					ModuleID: ModuleID,
					Code:     CodeSomeError,
					Message:  "some error",
				}
			},
			wantErrModID: ModuleID,
			wantErrCode:  CodeSomeError,
		},
		{
			name:     "non-elrpc err",
			moduleID: ModuleID,
			methodID: MethodID_HostAPI_Ping,
			req: &message.String{
				Value: "Ping",
			},
			pingImpl: func(arg *message.String) (*message.String, error) {
				return nil, errors.New("abnormal error")
			},
			wantErrModID: builtin.ModuleID,
			wantErrCode:  builtin.CodeInternal,
		},
		{
			name:     "no such method",
			moduleID: ModuleID,
			methodID: MethodID_HostAPI_NoSuchMethod,
			req: &message.String{
				Value: "Ping",
			},
			pingImpl: func(arg *message.String) (*message.String, error) {
				return &message.String{Value: "Pong"}, nil
			},
			wantErrModID: builtin.ModuleID,
			wantErrCode:  builtin.CodeUnimplemented,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hostImpl := &hostAPIImpl{}
			hostImpl.pingImpl = tt.pingImpl

			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
			mod := elrpctest.NewTestModule(t)
			defer mod.Close()

			instance := runtime.NewInstance(logger, mod)
			UseHostAPI(instance, hostImpl)
			err := instance.Start()
			if err != nil {
				t.Fatal(err)
			}

			s := mod.GuestStream()

			// Call HostAPI.Ping
			enc := message.NewEncoder()
			err = enc.EncodeUint32(tt.moduleID)
			if err != nil {
				t.Fatal(err)
			}
			err = enc.EncodeUint32(tt.methodID)
			if err != nil {
				t.Fatal(err)
			}
			err = tt.req.MarshalELRPC(enc)
			if err != nil {
				t.Fatal(err)
			}
			buf := enc.Buffer()
			lenBuf, err := message.AppendLength(nil, len(buf))
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.Write(lenBuf)
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.Write(buf)
			if err != nil {
				t.Fatal(err)
			}

			// Receive the response
			lenBuf = make([]byte, message.LengthSize)
			_, err = io.ReadFull(s, lenBuf)
			if err != nil {
				t.Fatal(err)
			}
			respLen, err := message.DecodeLength(lenBuf)
			if err != nil {
				t.Fatal(err)
			}
			buf = make([]byte, respLen)
			_, err = io.ReadFull(s, buf)
			if err != nil {
				t.Fatal(err)
			}

			dec := message.NewDecoder(buf)
			got := &Result{}
			err = got.UnmarshalELRPC(dec)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantOk != nil {
				want := &Result{IsOk: true, Ok: tt.wantOk}
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			} else {
				if got.IsOk {
					t.Errorf("want error but got %#v", got)
					return
				}
				elrpctest.AssertError(t, got.Err, tt.wantErrModID, tt.wantErrCode)
			}
		})
	}
}
