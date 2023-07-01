package runtime_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/apibuilder"
	"github.com/genkami/elsi/elrpc/elrpctest"
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/runtime"
	"github.com/genkami/elsi/elrpc/types"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

const (
	ModuleID = 0x0000_ffff // only used by test package

	MethodID_HostAPI_Ping         = 0x0000_1234
	MethodID_HostAPI_NoSuchMethod = 0x0000_ffff

	MethodID_GuestAPI_Ping = 0x0009_8765

	CodeSomeError = 0x0000_abcd
)

type HostAPI interface {
	Ping(*message.String) (*message.String, error)
}

func ImportHostAPI(rt types.Runtime, h HostAPI) {
	rt.Use(ModuleID, MethodID_HostAPI_Ping, apibuilder.TypedHandler1[*message.String, *message.String](h.Ping))
}

type hostAPIImpl struct {
	pingImpl func(*message.String) (*message.String, error)
}

func (h *hostAPIImpl) Ping(args *message.String) (*message.String, error) {
	return h.pingImpl(args)
}

type GuestAPI interface {
	Ping(*message.String) (*message.String, error)
}

type guestAPIImpl struct {
	pingImpl *apibuilder.MethodCaller1[*message.String, *message.String]
}

func ExportGuestAPI(rt types.Runtime) GuestAPI {
	return &guestAPIImpl{
		pingImpl: apibuilder.NewMethodCaller1[*message.String, *message.String](rt, ModuleID, MethodID_GuestAPI_Ping),
	}
}

func (g *guestAPIImpl) Ping(args *message.String) (*message.String, error) {
	return g.pingImpl.Call(args)
}

func TestInstance_callHostAPI(t *testing.T) {
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
			guest := elrpctest.NewTestGuest(t)
			defer guest.Close()

			rt := runtime.NewRuntime(logger, guest)
			ImportHostAPI(rt, hostImpl)
			err := rt.Start()
			if err != nil {
				t.Fatal(err)
			}

			s := guest.GuestStream()

			// Call HostAPI.Ping
			respDec := callHostAPI(t, s, tt.moduleID, tt.methodID, tt.req)
			got := &Result{}
			err = got.UnmarshalELRPC(respDec)
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

func TestInstance_callGuestAPI(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	guest := elrpctest.NewTestGuest(t)
	defer guest.Close()

	rt := runtime.NewRuntime(logger, guest)
	guestAPI := ExportGuestAPI(rt)
	err := rt.Start()
	if err != nil {
		t.Fatal(err)
	}

	startPoll := make(chan struct{})
	var eg errgroup.Group
	eg.Go(func() error {
		s := guest.GuestStream()

		<-startPoll

		// Call builtin.Exporter.PollMethodCall
		type PollResult = message.Result[*builtin.MethodCall, *message.Error]
		var respDec *message.Decoder
		var mCall *builtin.MethodCall
		for {
			respDec = callHostAPI(
				t, s,
				builtin.ModuleID, builtin.MethodID_Exporter_PollMethodCall)
			pollResult := &PollResult{}
			err = pollResult.UnmarshalELRPC(respDec)
			if err != nil {
				t.Fatal(err)
			}
			if !pollResult.IsOk {
				err := pollResult.Err
				if err.ModuleID == builtin.ModuleID && err.Code == builtin.CodeNotFound {
					continue
				}
				t.Fatalf("want ok but got %#v", pollResult)
			}
			mCall = pollResult.Ok
			if mCall.ModuleID != ModuleID {
				t.Fatalf("want %X but got %X", ModuleID, mCall.ModuleID)
			}
			if mCall.MethodID != MethodID_GuestAPI_Ping {
				t.Fatalf("want %X but got %X", MethodID_GuestAPI_Ping, mCall.MethodID)
			}
			argDec := message.NewDecoder(mCall.Args.Raw)
			arg := &message.String{}
			err = arg.UnmarshalELRPC(argDec)
			if err != nil {
				t.Fatal(err)
			}
			if arg.Value != "Ping" {
				t.Fatalf("want Ping but got %s", arg.Value)
			}
			break
		}

		// Call builtin.Exporter.SendResult
		type SendResultResult = message.Result[message.Void, *message.Error]
		rvEnc := message.NewEncoder()
		err = rvEnc.EncodeString("Pong")
		if err != nil {
			t.Fatal(err)
		}
		respDec = callHostAPI(
			t, s,
			builtin.ModuleID, builtin.MethodID_Exporter_SendResult,
			&builtin.MethodResult{
				CallID: mCall.CallID,
				RetVal: &message.Result[*message.Any, *message.Error]{
					IsOk: true,
					Ok:   &message.Any{Raw: rvEnc.Buffer()},
				},
			},
		)
		sendResultResult := &SendResultResult{}
		err = sendResultResult.UnmarshalELRPC(respDec)
		if err != nil {
			t.Fatal(err)
		}
		if !sendResultResult.IsOk {
			t.Fatalf("want ok but got %#v", sendResultResult)
		}

		return nil
	})

	eg.Go(func() error {
		got, err := guestAPI.Ping(&message.String{Value: "Ping"})
		if err != nil {
			t.Fatal(err)
		}
		if got.Value != "Pong" {
			t.Errorf("want Pong but got %s", got.Value)
		}
		return nil
	})

	close(startPoll)

	err = eg.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func callHostAPI(t *testing.T, s runtime.Stream, modID, methodID uint32, args ...message.Message) *message.Decoder {
	enc := message.NewEncoder()
	err := enc.EncodeUint32(modID)
	if err != nil {
		t.Fatal(err)
	}
	err = enc.EncodeUint32(methodID)
	if err != nil {
		t.Fatal(err)
	}
	for _, arg := range args {
		err = arg.MarshalELRPC(enc)
		if err != nil {
			t.Fatal(err)
		}
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

	return message.NewDecoder(buf)
}
