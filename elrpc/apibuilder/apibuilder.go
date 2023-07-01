package apibuilder

import (
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

type HostHandler0[R message.Message] func() (R, error)

func (h HostHandler0[R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	return h()
}

type HostHandler1[T1, R message.Message] func(T1) (R, error)

func (h HostHandler1[T1, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1))
}

type HostHandler2[T1, T2, R message.Message] func(T1, T2) (R, error)

func (h HostHandler2[T1, T2, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := message.NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2))
}

type HostHandler3[T1, T2, T3, R message.Message] func(T1, T2, T3) (R, error)

func (h HostHandler3[T1, T2, T3, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := message.NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := message.NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3))
}

type HostHandler4[T1, T2, T3, T4, R message.Message] func(T1, T2, T3, T4) (R, error)

func (h HostHandler4[T1, T2, T3, T4, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := message.NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := message.NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x4 := message.NewMessage[T4]()
	err = x4.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3), x4.(T4))
}

type HostHandler5[T1, T2, T3, T4, T5, R message.Message] func(T1, T2, T3, T4, T5) (R, error)

func (h HostHandler5[T1, T2, T3, T4, T5, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x2 := message.NewMessage[T2]()
	err = x2.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x3 := message.NewMessage[T3]()
	err = x3.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x4 := message.NewMessage[T4]()
	err = x4.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	x5 := message.NewMessage[T5]()
	err = x5.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1), x2.(T2), x3.(T3), x4.(T4), x5.(T5))
}

type delegator struct {
	rt       types.Runtime
	moduleID uint32
	methodID uint32
}

type GuestDelegator0[R message.Message] struct {
	delegator
}

func NewGuestDelegator0[R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator0[R] {
	return &GuestDelegator0[R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator0[R]) Call() (R, error) {
	var zero R
	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type GuestDelegator1[T1, R message.Message] struct {
	delegator
}

func NewGuestDelegator1[T1, R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator1[T1, R] {
	return &GuestDelegator1[T1, R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator1[T1, R]) Call(x1 T1) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type GuestDelegator2[T1, T2, R message.Message] struct {
	delegator
}

func NewGuestDelegator2[T1, T2, R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator2[T1, T2, R] {
	return &GuestDelegator2[T1, T2, R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator2[T1, T2, R]) Call(x1 T1, x2 T2) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type GuestDelegator3[T1, T2, T3, R message.Message] struct {
	delegator
}

func NewGuestDelegator3[T1, T2, T3, R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator3[T1, T2, T3, R] {
	return &GuestDelegator3[T1, T2, T3, R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator3[T1, T2, T3, R]) Call(x1 T1, x2 T2, x3 T3) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type GuestDelegator4[T1, T2, T3, T4, R message.Message] struct {
	delegator
}

func NewGuestDelegator4[T1, T2, T3, T4, R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator4[T1, T2, T3, T4, R] {
	return &GuestDelegator4[T1, T2, T3, T4, R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator4[T1, T2, T3, T4, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x4.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}

type GuestDelegator5[T1, T2, T3, T4, T5, R message.Message] struct {
	delegator
}

func NewGuestDelegator5[T1, T2, T3, T4, T5, R message.Message](rt types.Runtime, moduleID, methodID uint32) *GuestDelegator5[T1, T2, T3, T4, T5, R] {
	return &GuestDelegator5[T1, T2, T3, T4, T5, R]{
		delegator: delegator{
			rt:       rt,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *GuestDelegator5[T1, T2, T3, T4, T5, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4, x5 T5) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x2.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x3.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x4.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}
	err = x5.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.rt.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
	if err != nil {
		return zero, err
	}

	dec := message.NewDecoder(rawResp.Raw)
	resp := message.NewMessage[R]()
	err = resp.UnmarshalELRPC(dec)
	if err != nil {
		return zero, err
	}
	return resp.(R), nil
}
