package expimpl

import (
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elsi/api/exp"
	"golang.org/x/exp/slog"
)

type HTTP struct {
	logger *slog.Logger
	hs     *HandleSet
}

type httpListener struct {
	logger  *slog.Logger
	hs      *HandleSet
	lis     net.Listener
	waiters httpWaiterSet
}

var (
	_ io.Closer    = (*httpListener)(nil)
	_ http.Handler = (*httpListener)(nil)
)

type httpWaiterSet struct {
	mu    sync.Mutex
	next  uint64
	all   map[uint64]*httpWaiter
	queue []*exp.ServerRequest
}

type httpWaiter struct {
	respHeaderCh    chan *exp.ServerResponseHeader
	respHandleCh    chan *exp.Handle
	respBodyCloseCh chan struct{}
}

var _ exp.HTTP = (*HTTP)(nil)

func NewHTTP(logger *slog.Logger, hs *HandleSet) *HTTP {
	return &HTTP{
		logger: logger,
		hs:     hs,
	}
}

func (h *HTTP) Listen(addrAndPort *message.String) (*exp.Handle, error) {
	// TODO: restrict access
	lis, err := net.Listen("tcp", addrAndPort.Value)
	if err != nil {
		// TODO: convert to ELRPC error
		return nil, err
	}
	logger := h.logger.With(slog.String("addr", addrAndPort.Value))
	listener := &httpListener{
		logger: logger,
		hs:     h.hs,
		lis:    lis,
		waiters: httpWaiterSet{
			all: make(map[uint64]*httpWaiter),
		},
	}
	go func() {
		err := http.Serve(lis, listener)
		if err != nil && err != http.ErrServerClosed {
			logger.Error("server terminated unexpectedly",
				slog.String("error", err.Error()))
		}
	}()

	lisH := h.hs.Register(listener)
	return &exp.Handle{ID: lisH}, nil
}

func (h *HTTP) PollRequest(handle *exp.Handle) (*exp.ServerRequest, error) {
	lisAny, ok := h.hs.Get(handle.ID)
	if !ok {
		return nil, errNoSuchHandle
	}
	lis, ok := lisAny.(*httpListener)
	if !ok {
		return nil, errUnsupported
	}
	return lis.pollRequest()
}

func (h *HTTP) SendResponseHeader(handle *exp.Handle, reqID *message.Uint64, header *exp.ServerResponseHeader) (*exp.Handle, error) {
	lisAny, ok := h.hs.Get(handle.ID)
	if !ok {
		return nil, errNoSuchHandle
	}
	lis, ok := lisAny.(*httpListener)
	if !ok {
		return nil, errUnsupported
	}
	return lis.sendResponseHeader(reqID, header)
}

func (lis *httpListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqHandle := lis.hs.Register(r.Body)
	defer func() {
		rInstance, ok := lis.hs.Remove(reqHandle)
		if ok {
			closer, ok := rInstance.(io.Closer)
			if ok {
				_ = closer.Close()
			}
		}
	}()

	req := &exp.ServerRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   &exp.Handle{ID: reqHandle},
	}
	waiter := lis.enqueue(req)

	// meanwhile:
	// * guest calls HTTP.PollRequest
	// * guest calls Stream.Read to reqHandle
	// * guest calls HTTP.SendResponseHeader

	respHeader := <-waiter.respHeaderCh
	w.WriteHeader(int(respHeader.Status))

	respHandle := lis.hs.Register(&httpResponseWriter{
		w:      w,
		waiter: waiter,
	})
	waiter.respHandleCh <- &exp.Handle{ID: respHandle}

	// meanwhile:
	// * guest calls Stream.Write to respHandle
	// * guest calls Stream.Close to respHandle

	<-waiter.respBodyCloseCh
}

func (lis *httpListener) enqueue(req *exp.ServerRequest) *httpWaiter {
	w := &httpWaiter{
		respHeaderCh:    make(chan *exp.ServerResponseHeader, 1),
		respHandleCh:    make(chan *exp.Handle, 1),
		respBodyCloseCh: make(chan struct{}, 1),
	}
	lis.waiters.mu.Lock()
	defer lis.waiters.mu.Unlock()
	lis.waiters.next++
	lis.waiters.all[lis.waiters.next] = w
	lis.waiters.queue = append(lis.waiters.queue, req)
	return w
}

func (lis *httpListener) pollRequest() (*exp.ServerRequest, error) {
	lis.waiters.mu.Lock()
	defer lis.waiters.mu.Unlock()
	if len(lis.waiters.queue) == 0 {
		return nil, errNoRequest
	}
	req := lis.waiters.queue[0]
	lis.waiters.queue = lis.waiters.queue[1:]
	return req, nil
}

func (list *httpListener) sendResponseHeader(reqID *message.Uint64, header *exp.ServerResponseHeader) (*exp.Handle, error) {
	list.waiters.mu.Lock()
	w, ok := list.waiters.all[reqID.Value]
	list.waiters.mu.Unlock()
	if !ok {
		return nil, errNoSuchHandle
	}
	w.respHeaderCh <- header
	respHandle := <-w.respHandleCh
	return respHandle, nil
}

func (l *httpListener) Close() error {
	return l.lis.Close()
}

type httpResponseWriter struct {
	w      http.ResponseWriter
	waiter *httpWaiter
}

var _ io.WriteCloser = (*httpResponseWriter)(nil)

func (w *httpResponseWriter) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *httpResponseWriter) Close() error {
	w.waiter.respBodyCloseCh <- struct{}{}
	return nil
}
