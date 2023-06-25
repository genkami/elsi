package runtime

import (
	"fmt"
	"os"
	"sync"

	"github.com/genkami/elsi/elrpc/api/builtin"
	"github.com/genkami/elsi/elrpc/message"
)

type callResult struct {
	retVal *message.Result[*message.Any, *message.Error]
}

type exporterImpl struct {
	mu        sync.Mutex
	waiters   map[uint64]chan<- callResult
	callQueue []*builtin.MethodCall
	next      uint64
}

var _ builtin.Exporter = &exporterImpl{}

func newExporter() *exporterImpl {
	return &exporterImpl{
		waiters: make(map[uint64]chan<- callResult),
	}
}

func (e *exporterImpl) callAsync(call *builtin.MethodCall) <-chan callResult {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := make(chan callResult, 1)
	id := e.next
	e.next++
	call.CallID = id
	e.callQueue = append(e.callQueue, call)
	e.waiters[id] = ch
	return ch
}

func (e *exporterImpl) PollMethodCall() (*builtin.MethodCall, error) {
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

func (e *exporterImpl) SendResult(m *builtin.MethodResult) (*message.Void, error) {
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
	ch <- callResult{m.RetVal}
	return &message.Void{}, nil
}

func fullID(moduleID, methodID uint32) uint64 {
	return uint64(moduleID)<<32 | uint64(methodID)
}
