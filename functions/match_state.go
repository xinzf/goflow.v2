package functions

import (
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
)

type MatchState struct {
}

func (this *MatchState) GetName() string {
	return "matchState"
}

func (this *MatchState) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {
	return nil, nil
}
