package elrpctest

import (
	"os"
	"testing"

	"github.com/genkami/elsi/elrpc/runtime"
)

type TestModule struct {
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

var _ runtime.Module = (*TestModule)(nil)

func NewTestModule(t *testing.T) *TestModule {
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
	return &TestModule{
		hostEnd:  &pipeStream{r: hostR, w: hostW},
		guestEnd: &pipeStream{r: guestR, w: guestW},
	}
}

func (m *TestModule) Stream() runtime.Stream {
	return m.hostEnd
}

func (m *TestModule) Start() error {
	return nil
}

func (m *TestModule) Wait() error {
	return nil
}

func (m *TestModule) GuestStream() runtime.Stream {
	return m.guestEnd
}

func (m *TestModule) Close() {
	m.hostEnd.Close()
	m.guestEnd.Close()
}
