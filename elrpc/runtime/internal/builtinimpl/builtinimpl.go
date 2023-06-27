package builtinimpl

import (
	"fmt"
	"os"
	"sync"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
)

type CallResult struct {
	RetVal *message.Result[*message.Any, *message.Error]
}

type Exporter struct {
	mu        sync.Mutex
	waiters   map[uint64]chan<- CallResult
	callQueue []*builtin.MethodCall
	next      uint64
}

var _ builtin.Exporter = &Exporter{}

func NewExporter() *Exporter {
	return &Exporter{
		waiters: make(map[uint64]chan<- CallResult),
	}
}

func (e *Exporter) CallAsync(call *builtin.MethodCall) <-chan CallResult {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := make(chan CallResult, 1)
	id := e.next
	e.next++
	call.CallID = id
	e.callQueue = append(e.callQueue, call)
	e.waiters[id] = ch
	return ch
}

func (e *Exporter) PollMethodCall() (*builtin.MethodCall, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.callQueue) == 0 {
		return nil, &message.Error{
			Code:    builtin.CodeNotFound,
			Message: "no method call",
		}
	}
	call := e.callQueue[0]
	e.callQueue = e.callQueue[1:]
	return call, nil
}

func (e *Exporter) SendResult(m *builtin.MethodResult) (*message.Void, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.waiters[m.CallID]
	if !ok {
		// TODO: slog
		fmt.Fprintf(os.Stderr, "no such call: %X\n", m.CallID)
		return nil, &message.Error{
			Code:    builtin.CodeNotFound,
			Message: "no such method call",
		}
	}
	ch <- CallResult{m.RetVal}
	return &message.Void{}, nil
}
