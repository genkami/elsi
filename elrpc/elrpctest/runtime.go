package elrpctest

import (
	"os"
	"testing"

	"github.com/genkami/elsi/elrpc/runtime"
)

type TestGuest struct {
	hostEnd  *pipeStream
	guestEnd *pipeStream
}

type pipeStream struct {
	r *os.File
	w *os.File
}

func (s *pipeStream) Read(p []byte) (int, error) {
	return s.r.Read(p)
}

func (s *pipeStream) Write(p []byte) (int, error) {
	return s.w.Write(p)
}

func (s *pipeStream) Close() error {
	err := s.r.Close()
	if err != nil {
		_ = s.w.Close()
		return err
	}
	return s.w.Close()
}

var _ runtime.Guest = (*TestGuest)(nil)

func NewTestGuest(t *testing.T) *TestGuest {
	hostR, guestW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create a pipe: %s", err.Error())
	}
	guestR, hostW, err := os.Pipe()
	if err != nil {
		hostR.Close()
		guestW.Close()
		t.Fatalf("failed to create a pipe: %s", err.Error())
	}
	return &TestGuest{
		hostEnd:  &pipeStream{r: hostR, w: hostW},
		guestEnd: &pipeStream{r: guestR, w: guestW},
	}
}

func (g *TestGuest) Stream() runtime.Stream {
	return g.hostEnd
}

func (g *TestGuest) Start() error {
	return nil
}

func (g *TestGuest) Wait() error {
	return nil
}

func (g *TestGuest) GuestStream() runtime.Stream {
	return g.guestEnd
}

func (g *TestGuest) Close() {
	g.hostEnd.Close()
	g.guestEnd.Close()
}
