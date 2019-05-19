package spi

import (
	"github.com/xinzf/goflow.v2/tools"
)

type Function interface {
	GetName() string
	Eval(store Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error)
}
