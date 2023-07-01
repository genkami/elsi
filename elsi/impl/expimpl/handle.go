package expimpl

import "sync"

type HandleSet struct {
	mu    sync.Mutex
	next  uint64
	items map[uint64]any
}

func NewHandleSet() *HandleSet {
	return &HandleSet{
		items: make(map[uint64]any),
	}
}

func (hs *HandleSet) Register(v any) uint64 {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.next++
	hs.items[hs.next] = v
	return hs.next
}

func (hs *HandleSet) Get(handle uint64) (any, bool) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	v, ok := hs.items[handle]
	return v, ok
}

func (hs *HandleSet) Remove(handle uint64) (any, bool) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	v, ok := hs.items[handle]
	delete(hs.items, handle)
	return v, ok
}
