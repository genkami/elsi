// TODO: move this to somewhere else
package helpers

import (
	"github.com/genkami/elsi/elrpc/message"
	"github.com/genkami/elsi/elrpc/types"
)

type TypedHandler0[R message.Message] func() (R, error)

func (h TypedHandler0[R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	return h()
}

type TypedHandler1[T1, R message.Message] func(T1) (R, error)

func (h TypedHandler1[T1, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
	x1 := message.NewMessage[T1]()
	err := x1.UnmarshalELRPC(dec)
	if err != nil {
		return nil, err
	}

	return h(x1.(T1))
}

type TypedHandler2[T1, T2, R message.Message] func(T1, T2) (R, error)

func (h TypedHandler2[T1, T2, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
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

type TypedHandler3[T1, T2, T3, R message.Message] func(T1, T2, T3) (R, error)

func (h TypedHandler3[T1, T2, T3, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
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

type TypedHandler4[T1, T2, T3, T4, R message.Message] func(T1, T2, T3, T4) (R, error)

func (h TypedHandler4[T1, T2, T3, T4, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
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

type TypedHandler5[T1, T2, T3, T4, T5, R message.Message] func(T1, T2, T3, T4, T5) (R, error)

func (h TypedHandler5[T1, T2, T3, T4, T5, R]) HandleRequest(dec *message.Decoder) (message.Message, error) {
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

type methodCaller struct {
	instance types.Instance
	moduleID uint32
	methodID uint32
}

type MethodCaller0[R message.Message] struct {
	methodCaller
}

func NewMethodCaller0[R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller0[R] {
	return &MethodCaller0[R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller0[R]) Call() (R, error) {
	var zero R
	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{})
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

type MethodCaller1[T1, R message.Message] struct {
	methodCaller
}

func NewMethodCaller1[T1, R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller1[T1, R] {
	return &MethodCaller1[T1, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller1[T1, R]) Call(x1 T1) (R, error) {
	var zero R
	enc := message.NewEncoder()
	err := x1.MarshalELRPC(enc)
	if err != nil {
		return zero, err
	}

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
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

type MethodCaller2[T1, T2, R message.Message] struct {
	methodCaller
}

func NewMethodCaller2[T1, T2, R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller2[T1, T2, R] {
	return &MethodCaller2[T1, T2, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller2[T1, T2, R]) Call(x1 T1, x2 T2) (R, error) {
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

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
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

type MethodCaller3[T1, T2, T3, R message.Message] struct {
	methodCaller
}

func NewMethodCaller3[T1, T2, T3, R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller3[T1, T2, T3, R] {
	return &MethodCaller3[T1, T2, T3, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller3[T1, T2, T3, R]) Call(x1 T1, x2 T2, x3 T3) (R, error) {
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

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
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

type MethodCaller4[T1, T2, T3, T4, R message.Message] struct {
	methodCaller
}

func NewMethodCaller4[T1, T2, T3, T4, R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller4[T1, T2, T3, T4, R] {
	return &MethodCaller4[T1, T2, T3, T4, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller4[T1, T2, T3, T4, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4) (R, error) {
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

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
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

type MethodCaller5[T1, T2, T3, T4, T5, R message.Message] struct {
	methodCaller
}

func NewMethodCaller5[T1, T2, T3, T4, T5, R message.Message](instance types.Instance, moduleID, methodID uint32) *MethodCaller5[T1, T2, T3, T4, T5, R] {
	return &MethodCaller5[T1, T2, T3, T4, T5, R]{
		methodCaller: methodCaller{
			instance: instance,
			moduleID: moduleID,
			methodID: methodID,
		},
	}
}

func (c *MethodCaller5[T1, T2, T3, T4, T5, R]) Call(x1 T1, x2 T2, x3 T3, x4 T4, x5 T5) (R, error) {
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

	rawResp, err := c.instance.Call(c.moduleID, c.methodID, &message.Any{Raw: enc.Buffer()})
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
