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

func UseWorld(rt types.Runtime, imports *Imports) *Exports {
	ImportTODO(rt, imports.TODO)
	return &Exports{
		Greeter: ExportGreeter(rt),
	}
}
