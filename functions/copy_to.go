package functions

import (
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
)

type CopyTo struct {
}

func (this *CopyTo) GetName() string {
	return "copyTo"
}

func (this *CopyTo) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {
	return nil, nil
}
