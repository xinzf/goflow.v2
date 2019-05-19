package functions

import (
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
)

type CanJoin struct {
}

func (this *CanJoin) GetName() string {
	return "canJoin"
}

func (this *CanJoin) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {
	currentStep := transientVars.Get(tools.CurrentStep).GetData().(spi.Step)

	if args["step_id"].Int() == currentStep.GetStepId() {
		if args["state"].String() == currentStep.GetState() && args["action"].String() == currentStep.GetActionName() {
			return true, nil
		}
		return false, nil
	}

	trans, err := store.GetJoinTransitionsByPrevId(
		transientVars.Get(tools.Entry).GetData().(spi.Entry).GetEntryId(),
		args["step_id"].Int(),
	)
	if err != nil {
		return nil, err
	}

	if trans.GetState() == args["state"].String() && trans.GetActionName() == args["action"].String() {
		return true, nil
	}
	return false, nil
}
