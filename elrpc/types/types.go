package types

import (
	"github.com/genkami/elsi/elrpc/message"
)

type Runtime interface {
	Use(moduleID uint32, methodID uint32, handler HostHandler)
	Call(moduleID uint32, methodID uint32, args *message.Any) (*message.Any, error)
}

type HostHandler interface {
	HandleRequest(*message.Decoder) (message.Message, error)
}
