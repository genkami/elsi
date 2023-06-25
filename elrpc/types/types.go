package types

import (
	"github.com/genkami/elsi/elrpc/message"
)

type Instance interface {
	Use(moduleID uint32, methodID uint32, handler Handler)
	Call(moduleID uint32, methodID uint32, args *message.Any) (*message.Any, error)
}

type Handler interface {
	HandleRequest(*message.Decoder) (message.Message, error)
}
