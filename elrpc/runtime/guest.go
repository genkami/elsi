package runtime

import (
	"io"
	"os"
	"os/exec"
)

type Stream interface {
	io.ReadWriter
}

type pipeStream struct {
	io.Reader
	io.Writer
}

func NewPipeStream(in io.Reader, out io.Writer) Stream {
	return &pipeStream{in, out}
}

type Guest interface {
	Stream() Stream
	Start() error
	Wait() error
}

type ProcessGuest struct {
	cmd    *exec.Cmd
	stream Stream
}

var _ Guest = (*ProcessGuest)(nil)

func NewProcessGuest(name string, args ...string) *ProcessGuest {
	cmd := exec.Command(name, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr
	return &ProcessGuest{
		cmd:    cmd,
		stream: NewPipeStream(stdout, stdin),
	}
}

func (m *ProcessGuest) Stream() Stream {
	return m.stream
}

func (m *ProcessGuest) Start() error {
	return m.cmd.Start()
}

func (m *ProcessGuest) Wait() error {
	return m.cmd.Wait()
}
