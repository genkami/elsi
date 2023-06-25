package x

import (
	"github.com/genkami/elsi/elrpc"
)

type Exports struct {
	Greeter Greeter
}

func UseWorld(instance *elrpc.Instance, todo TODO) *Exports {
	ImportTODO(instance, todo)
	return &Exports{
		Greeter: ExportGreeter(instance),
	}
}
