package x

import (
	"github.com/genkami/elsi/elrpc/types"
)

type Imports struct {
	TODO TODO
}

type Exports struct {
	Greeter Greeter
}

func UseWorld(instance types.Instance, imports *Imports) *Exports {
	ImportTODO(instance, imports.TODO)
	return &Exports{
		Greeter: ExportGreeter(instance),
	}
}
