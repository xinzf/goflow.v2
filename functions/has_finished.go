package functions

import (
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
)

type HasFinished struct {
}

func (this *HasFinished) GetName() string {
	return "hasFinished"
}

func (this *HasFinished) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {
	return nil, nil
}
