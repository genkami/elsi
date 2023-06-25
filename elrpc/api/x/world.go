package x

import (
	"github.com/genkami/elsi/elrpc/types"
)

type Exports struct {
	Greeter Greeter
}

func UseWorld(instance types.Instance, todo TODO) *Exports {
	ImportTODO(instance, todo)
	return &Exports{
		Greeter: ExportGreeter(instance),
	}
}
