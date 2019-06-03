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
	entry := transientVars.Get(tools.Entry).GetData().(spi.Entry)
	currentStep := transientVars.Get(tools.CurrentStep).GetData().(spi.Step)

	// 当前步骤是否就是要验证的步骤
	// 如果是，就直接验证
	if args["step_id"].Int() == currentStep.GetStepId() {
		if args["state"].String() == currentStep.GetState() && args["action"].String() == currentStep.GetActionName() {
			return true, nil
		}
		return false, nil
	}

	// 如果当前步骤不是要验证的步骤，就从数据库中调出来并验证
	// 假设当前步骤依然是在 currentStep 中，说明该步骤还未结束，直接验证为假，并退出
	_, found, err := store.FindCurrentStep(
		entry.GetEntryId(),
		args["step_id"].Int(),
	)
	if err != nil {
		return false, err
	}
	if found {
		return false, nil
	}

	// 如果要验证的步骤在历史记录中
	history, err := store.FindMostRecentHistory(
		entry.GetEntryId(),
		args["step_id"].Int(),
	)
	if err != nil {
		return false, err
	}

	// 如果历史记录中的记录不满足条件
	if history.GetState() != args["state"].String() || history.GetActionName() != args["action"].String() {
		return false, nil
	}

	// 这个已执行完成的步骤，他有没有前置未完成的任务，如果有，也不能通过汇集条件
	prevIds := history.GetPrevIds()
	for _, id := range prevIds {
		_, found, err := store.FindCurrentStep(
			entry.GetEntryId(),
			id,
		)
		if err != nil {
			return false, err
		}
		if found {
			return false, nil
		}
	}

	return true, nil

	// 如果要验证的步骤已经不再 currentStep 中了，那么则意味着两种可能
	// 第一种可能，这个步骤已经执行过了，现在在 join_transition 中等待
	// 第二种可能，这个步骤还未生成（可能还有其他前置认为未执行完）
	// 所以归根结底，就是从 join_transition 中查询要验证的目标任务
	//trans, err := store.GetJoinTransitionsByPrevId(
	//	transientVars.Get(tools.Entry).GetData().(spi.Entry).GetEntryId(),
	//	args["step_id"].Int(),
	//)
	//if err != nil {
	//	return nil, err
	//}

	//if trans.GetState() == args["state"].String() && trans.GetActionName() == args["action"].String() {
	//	return true, nil
	//}
	//return false, nil
}
