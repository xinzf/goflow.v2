package goflow

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xinzf/goflow.v2/definition"
	"github.com/xinzf/goflow.v2/enums"
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
	"math/rand"
	"time"
)

func NewTransition(set *PropertSet, vars *tools.TransientVars, prevStep, currentStep spi.Step) *Transition {
	vars.Put(tools.PrevStep, prevStep)
	vars.Put(tools.CurrentStep, currentStep)
	return &Transition{
		properSet: set,
		vars:      vars,
	}
}

type Transition struct {
	properSet *PropertSet
	vars      *tools.TransientVars
}

// @todo 这个没用到，应该在跳出 step 的时候调用
// 当前步骤退出时，触发的事件
func (this *Transition) exit(step definition.Step, action definition.Action, result definition.Result, remark string) error {
	currentStep := this.vars.Get(tools.CurrentStep).GetData().(spi.Step)
	currentStep.SetState(result.ExitStatus)
	currentStep.SetFinishDate(time.Now())
	currentStep.SetCaller(this.vars.Get(tools.Caller).Int())
	if err := this.properSet.GetStore().MoveHistory(currentStep); err != nil {
		return err
	}
	for _, fun := range step.PostFunctions.Functions {
		if _, err := this.triggerFunction(fun); err != nil {
			return err
		}
	}
	return nil
}

// 当前步骤进入时，触发的事件
func (this *Transition) Enter() error {
	workflow := this.vars.Get(tools.Workflow).GetData().(definition.Workflow)
	currentStep := this.vars.Get(tools.CurrentStep).GetData().(spi.Step)
	step, _ := workflow.GetStep(currentStep.GetStepId())

	// 验证限制条件
	if flag, err := this.evalConditions(step.Restrict.Conditions); err != nil {
		return nil
	} else if !flag {
		return errors.New(step.Restrict.GetMessage())
	}

	// 执行步骤前置事件
	for _, fun := range step.PreFunctions.Functions {
		if _, err := this.triggerFunction(fun); err != nil {
			return err
		}
	}

	// 没有可执行操作，流程结束
	if len(step.Actions.Actions) == 0 {
		entry := this.vars.Get(tools.Entry).GetData().(spi.Entry)
		err := this.properSet.store.DeleteAllCurrentStep(entry.GetEntryId())
		if err != nil {
			return err
		}

		entry.SetEndTime(time.Now())
		entry.SetState(enums.COMPLETED)
		return this.properSet.GetStore().UpdateEntry(entry)
	}

	// 找出当前步骤中是否存在自动执行的 action
	// 如果有自动执行的 action，则自动触发
	autoActions := step.GetAutoActions()
	if len(autoActions) > 0 {
		for _, a := range autoActions {
			// 这里必须要 new 新的 transition，否则会存在数据污染的问题，每一个当前步骤都有一个独属于自己的 transition和临时变量池
			trans := NewTransition(
				this.properSet,
				this.vars,
				this.vars.Get(tools.PrevStep).GetData().(spi.Step),
				currentStep,
			)
			if err := trans.DoAction(a.Name, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Transition) beforeActions(actions definition.Actions) error {
	flag, err := this.evalConditions(actions.Restrict.Conditions)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New(actions.Restrict.GetMessage())
	}
	return nil
}

func (this *Transition) beforeAction(action definition.Action) error {
	flag, err := this.evalConditions(action.Restrict.Conditions)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New(action.Restrict.GetMessage())
	}

	for _, fun := range action.PreFunctions.Functions {
		if _, err = this.triggerFunction(fun); err != nil {
			return err
		}
	}

	return nil
}

func (this *Transition) afterAction(action definition.Action) error {
	for _, fun := range action.PostFunctions.Functions {
		if _, err := this.triggerFunction(fun); err != nil {
			return err
		}
	}
	return nil
}

func (this *Transition) DoAction(actionName, remark string) (err error) {
	workflow := this.vars.Get(tools.Workflow).GetData().(definition.Workflow)
	currentStep := this.vars.Get(tools.CurrentStep).GetData().(spi.Step)
	entry := this.vars.Get(tools.Entry).GetData().(spi.Entry)
	step, found := workflow.GetStep(currentStep.GetStepId())
	if !found {
		return fmt.Errorf("Not found step with step id: %d", currentStep.GetStepId())
	}
	action, found := step.GetAction(actionName)
	if !found {
		return fmt.Errorf("Not found action with: %s", actionName)
	}

	// beforeAction
	if err = this.beforeActions(step.Actions); err != nil {
		return err
	}
	if err = this.beforeAction(action); err != nil {
		return err
	}

	results, err := this.getResult(action.Results)
	if err != nil {
		return err
	}

	result := results[0]

	if err = this.afterAction(action); err != nil {
		return err
	}

	if err = this.transition(workflow, currentStep, entry, step, result, actionName, remark); err != nil {
		return err
	}

	if len(workflow.Global.PostFunctions.Functions) > 0 {
		for _, fn := range workflow.Global.PostFunctions.Functions {
			if _, err = this.triggerFunction(fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Transition) transition(workflow definition.Workflow, currentStep spi.Step, entry spi.Entry, step definition.Step, result definition.Result, actionName, remark string) (err error) {

	action, found := step.GetAction(actionName)
	if !found {
		return fmt.Errorf("Not found action with: %s", actionName)
	}

	removeReapetStepIds := func(stepIds []int) []int {
		newSlice := make([]int, 0)
		maps := make(map[int]int)

		for _, id := range stepIds {
			if _, ok := maps[id]; ok {
				continue
			}
			newSlice = append(newSlice, id)
			maps[id] = 0
		}

		return newSlice
	}

	nextStepId := this.getNextStepId(result.Step)

	currentStep.SetAction(action)
	currentStep.SetRemark(remark)

	if this.isExit(result) {
		if err = this.exit(step, action, result, remark); err != nil {
			return err
		}
	}

	// 当前步骤循环，没有跳出当前步骤
	if nextStepId == -1 {
		if result.Status != "" {
			currentStep.SetState(result.Status)
			if err = this.properSet.GetStore().UpdateCurrentStep(currentStep); err != nil {
				return err
			}
		}
		return nil
	}

	// 当前步骤跳出，但是整个流程结束，不会产生下一个步骤
	if nextStepId == -2 {
		if err = this.properSet.GetStore().DeleteAllCurrentStep(entry.GetEntryId()); err != nil {
			return err
		}

		entry.SetState(enums.KILLED)
		entry.SetEndTime(time.Now())
		if err = this.properSet.GetStore().UpdateEntry(entry); err != nil {
			return err
		}

		return nil
	}

	// 顺序进入下一个步骤，回退只会出现在顺序的情况下
	if nextStepId > 0 {

		nextStep, found := workflow.GetStep(nextStepId)
		if !found {
			return fmt.Errorf("Not found step with step_id: %d", nextStepId)
		}

		owner, err := this.getOwner(result.Owners.Owners)
		if err != nil {
			return err
		}

		// action 是一个回退操作
		if this.isPrevId(nextStepId, currentStep.GetPrevIds()) {
			if err := this.rollback(nextStep, owner, result.GetDueTime(), result.Status); err != nil {
				return err
			}
			return nil
		}

		prevIds := removeReapetStepIds(append(currentStep.GetPrevIds(), currentStep.GetStepId()))
		newCurrentStep, err := this.properSet.GetStore().CreateCurrentStep(
			entry.GetEntryId(),
			workflow,
			nextStep,
			owner,
			result.GetDueTime(),
			result.Status,
			prevIds,
		)
		if err != nil {
			return err
		}

		this.vars.Put(tools.NewSteps, []int{newCurrentStep.GetStepId()})

		trans := NewTransition(
			this.properSet,
			this.vars,
			currentStep,
			newCurrentStep,
		)
		if err = trans.Enter(); err != nil {
			return err
		}

		return nil
	}

	if result.Join > 0 {
		// 获取 join 中的 result
		join, _ := workflow.GetJoin(result.Join)
		if len(join.Results.Results) == 0 {
			return fmt.Errorf("Join's results definition is wrong", result.Join)
		}

		// 从 result 中获取要 join 的目标步骤ID
		result = join.Results.Results[0]
		nextStepId = this.getNextStepId(result.Step)

		// 验证 result 中的条件
		if verify, err := this.evalConditions(result.Conditions); err != nil {
			return err
		} else if !verify {
			//t, err := this.properSet.GetStore().GetJoinTransitionsByPrevId(entry.GetEntryId(), currentStep.GetStepId())
			//if err != nil {
			//	return err
			//}
			//if t.GetPrevId() != currentStep.GetStepId() {
			//	err = this.properSet.GetStore().CreateJoinTransition(entry.GetEntryId(), currentStep, nextStepId)
			//	return err
			//}
			return nil
		}

		nextStep, found := workflow.GetStep(nextStepId)
		if !found {
			return fmt.Errorf("Not found step with step_id: %d", nextStepId)
		}

		owner, err := this.getOwner(result.Owners.Owners)
		if err != nil {
			return err
		}

		//joinTransition, err := this.properSet.GetStore().GetJoinTransitionsByNextId(entry.GetEntryId(), nextStepId)
		//if err != nil {
		//	return err
		//}
		//
		//if err = this.properSet.GetStore().DeleteJoinTransition(entry.GetEntryId(), nextStepId); err != nil {
		//	return err
		//}
		//for _,r:=range join.Results.Results{
		//	r.Step
		//}

		prevIds := make([]int, 0)
		steps := workflow.GetSteps()
		for _, s := range steps {
			if s.IsTheJoin(join.Id) {
				stepHistory, err := this.properSet.GetStore().FindMostRecentHistory(entry.GetEntryId(), s.ID)
				if err != nil {
					return err
				}

				prevIds = append(prevIds, stepHistory.GetPrevIds()...)
				prevIds = append(prevIds, stepHistory.GetStepId())
			}
		}

		//for _, j := range joinTransition {
		//	joinStep, err := this.properSet.GetStore().FindMostRecentHistory(entry.GetEntryId(), j.GetPrevId())
		//	if err != nil {
		//		return err
		//	}

		//prevIds = append(prevIds, joinStep.GetPrevIds()...)
		//prevIds = append(prevIds, j.GetPrevId())
		//}
		prevIds = append(prevIds, currentStep.GetPrevIds()...)
		prevIds = append(prevIds, currentStep.GetStepId())
		prevIds = removeReapetStepIds(prevIds)

		newCurrentStep, err := this.properSet.GetStore().CreateCurrentStep(
			entry.GetEntryId(),
			workflow,
			nextStep,
			owner,
			result.GetDueTime(),
			result.Status,
			prevIds,
		)
		if err != nil {
			return err
		}

		this.vars.Put(tools.NewSteps, []int{newCurrentStep.GetStepId()})

		trans := NewTransition(
			this.properSet,
			this.vars,
			currentStep,
			newCurrentStep,
		)
		if err = trans.Enter(); err != nil {
			return err
		}
		return nil
	}

	if result.Split > 0 {
		split, _ := workflow.GetSplit(result.Split)
		results, err := this.getResult(split.Results)
		if err != nil {
			return err
		}

		newSteps := make([]int, 0)
		for _, result = range results {
			nextStep, found := workflow.GetStep(this.getNextStepId(result.Step))
			if !found {
				return fmt.Errorf("Not found step with step_id: %d", nextStepId)
			}
			// @todo
			owner, err := this.getOwner(result.Owners.Owners)
			if err != nil {
				return err
			}

			prevIds := removeReapetStepIds(append(currentStep.GetPrevIds(), currentStep.GetStepId()))
			newCurrentStep, err := this.properSet.GetStore().CreateCurrentStep(
				entry.GetEntryId(),
				workflow,
				nextStep,
				owner,
				result.GetDueTime(),
				result.Status,
				prevIds,
			)
			if err != nil {
				return err
			}
			newSteps = append(newSteps, newCurrentStep.GetStepId())

			trans := NewTransition(
				this.properSet,
				this.vars,
				currentStep,
				newCurrentStep,
			)
			if err = trans.Enter(); err != nil {
				return err
			}
		}
		this.vars.Put(tools.NewSteps, newSteps)

		return nil
	}

	return nil
}

func (this *Transition) triggerFunction(fun definition.Function) (interface{}, error) {
	function, found := this.properSet.GetFunction(fun.Name)
	if !found {
		return nil, fmt.Errorf("Not found the function with the name: %s", fun.Name)
	}

	return function.Eval(
		this.properSet.GetStore(),
		this.vars,
		fun.GetArgValues(),
	)
}

func (this *Transition) evalConditions(conditions definition.Conditions) (bool, error) {
	//return true, nil
	if len(conditions.Functions) == 0 {
		return true, nil
	}

	if conditions.Type.String() == "" {
		conditions.Type = enums.AND
	}

	var flag bool
	if conditions.Type == enums.AND {
		flag = true
		for _, f := range conditions.Functions {
			i, err := this.triggerFunction(f)
			if err != nil {
				return false, err
			}

			switch i.(type) {
			case bool:
				if i.(bool) != f.GetWant() {
					flag = false
					break
				}
			default:
				return false, fmt.Errorf("Function: %s's result is not boolean type", f.Name)
			}
		}
	} else if conditions.Type == enums.OR {
		flag = false
		for _, f := range conditions.Functions {
			i, err := this.triggerFunction(f)
			if err != nil {
				return false, err
			}

			switch i.(type) {
			case bool:
				if i.(bool) == f.GetWant() {
					flag = true
					break
				}
			default:
				return false, fmt.Errorf("Function: %s's result is not boolean type", f.Name)
			}
		}
	}

	return flag, nil
}

func (this *Transition) getResult(results definition.Results) ([]definition.Result, error) {
	finalResults := make([]definition.Result, 0)

	if len(results.Results) > 0 {
		for _, r := range results.Results {
			if verify, err := this.evalConditions(r.Conditions); err != nil {
				return nil, err
			} else if verify {
				finalResults = append(finalResults, r)
			}
		}
	}

	if len(finalResults) == 0 {
		finalResults = append(finalResults, results.Default...)
	}

	return finalResults, nil
}

func (this *Transition) getNextStepId(step string) int {
	var nextStep int
	val := tools.NewValue(step)
	if key, flag := val.ParseVariable(); flag {
		nextStep = this.vars.Get(key).Int()
	} else {
		nextStep = val.Int()
	}

	return nextStep
}

func (this *Transition) isExit(result definition.Result) bool {
	nextStepId := this.getNextStepId(result.Step)
	if result.Join > 0 || result.Split > 0 || nextStepId == -2 || nextStepId > 0 {
		return true
	}
	return false
}

func (this *Transition) getOwners(owners []definition.Owner) ([]int, error) {
	uids := make([]int, 0)
	for _, o := range owners {
		switch o.Type {
		case enums.Users:
			val := tools.NewValue(o.Props)
			uids = append(uids, val.IntSlice(",")...)
		case enums.Leader:
			user, err := this.properSet.GetStore().GetUser(this.vars.Get(tools.Caller).Int())
			if err != nil {
				return nil, err
			}

			leader, err := user.GetMyLeader()
			if err != nil {
				return nil, err
			}

			uids = append(uids, leader.GetId())
		case enums.Caller:
			uids = append(uids, this.vars.Get(tools.Caller).Int())
		case enums.Creator:
			uids = append(uids, this.vars.Get(tools.Entry).GetData().(spi.Entry).GetCreator())
		case enums.Roles:
			val := tools.NewValue(o.Props)
			users, err := this.properSet.GetStore().GetUsersByRoleIds(val.IntSlice(","))
			if err != nil {
				return nil, err
			}

			for _, u := range users {
				uids = append(uids, u.GetId())
			}
		case enums.Deps:
			val := tools.NewValue(o.Props)
			users, err := this.properSet.GetStore().GetUsersByRoleIds(val.IntSlice(","))
			if err != nil {
				return nil, err
			}

			for _, u := range users {
				uids = append(uids, u.GetId())
			}
		case enums.Variable:

			var (
				key  string
				flag bool
			)

			val := tools.NewValue(o.Props)
			if key, flag = val.ParseVariable(); !flag {
				return nil, fmt.Errorf("The key: %s is not correct variable format", o.Props)
			}

			uids = append(uids, this.vars.Get(key).Int())
		}
	}

	return uids, nil
}

func (this *Transition) getOwner(owners []definition.Owner) (int, error) {
	uids, err := this.getOwners(owners)
	if err != nil {
		return 0, err
	}

	if len(uids) == 0 {
		return 0, nil
	}

	if len(uids) == 1 {
		return uids[0], nil
	}

	// 随机取一个
	return uids[rand.Intn(len(uids)-1)], nil

	return uids[0], nil
}

func (this *Transition) isPrevId(stepId int, prevIds []int) bool {
	for _, id := range prevIds {
		if stepId == id {
			return true
		}
	}
	return false
}

// 回退操作
//func (this *Transition) rollback(step definition.Step, owner int, dueTime time.Time, state string) error {
//	logrus.Debugln("查询案件")
//	entry := this.vars.Get(tools.Entry).GetData().(spi.Entry)
//
//	logrus.Println("查询所有正在进行中的步骤")
//	allCurrentSteps, err := this.properSet.store.FindAllCurrentStep(entry.GetEntryId())
//	deleteIds := make([]int, 0)
//	//deleteIds := []int{
//	//	step.ID,
//	//}
//	for _, s := range allCurrentSteps {
//		// 如果要回退的步骤存在于这些步骤的前置中
//		// 那么就删除这些步骤
//		if this.isPrevId(step.ID, s.GetPrevIds()) {
//			deleteIds = append(deleteIds, s.GetStepId())
//		}
//	}
//	logrus.Debugln("得到了要删除的进行中的步骤：", deleteIds)
//	logrus.Debugln("删除进行中的步骤")
//	if err = this.properSet.store.DeleteCurrentSteps(deleteIds); err != nil {
//		return err
//	}
//
//	history, err := this.properSet.store.FindHistorySteps(entry.GetEntryId())
//	if err != nil {
//		return err
//	}
//
//	// @todo 这里要加注释，很难说清楚
//	for _, h := range history {
//		deleteIds = append(deleteIds, h.GetStepId())
//	}
//
//	logrus.Debugln("删除所有的join关系")
//	// 删除所有与当前步骤是兄弟步骤的 Join_transition 记录
//	// prev_ids 来源于 deleteIds，因为 deleteIds 代表的是所有要回退的步骤
//	// 那么理所当然的是这些要回退的步骤也要从 join_transition 中删除（因为没有必要等待，后续还会重新生成这些步骤）
//	deleteIds = append(deleteIds, step.ID)
//	if err = this.properSet.store.DeleteJoinTransitionByPrevIds(entry.GetEntryId(), deleteIds); err != nil {
//		return err
//	}
//
//	logrus.Debugln("移动步骤到历史记录")
//	// 从历史记录中找出之前的记录并把 prev_ids 置为其前置ID
//	historyStep, err := this.properSet.store.FindMostRecentHistory(entry.GetEntryId(), step.ID)
//	if err != nil {
//		return err
//	}
//
//	prevIds := historyStep.GetPrevIds()
//	logrus.Debugln("创建新的步骤")
//	// 创建新的步骤
//	newCurrentStep, err := this.properSet.store.CreateCurrentStep(
//		entry.GetEntryId(),
//		this.vars.Get(tools.Workflow).GetData().(definition.Workflow),
//		step,
//		owner,
//		dueTime,
//		state,
//		prevIds,
//	)
//
//	// 发起新路程并进入该步骤
//	trans := NewTransition(
//		this.properSet,
//		this.vars,
//		this.vars.Get(tools.CurrentStep).GetData().(spi.Step),
//		newCurrentStep,
//	)
//
//	logrus.Debugln("进入新创建的步骤")
//	if err = trans.Enter(); err != nil {
//		return err
//	}
//
//	return nil
//}

func (this *Transition) rollback(targetStep definition.Step, owner int, dueTime time.Time, state string) error {
	entry := this.vars.Get(tools.Entry).GetData().(spi.Entry)
	wf := this.vars.Get(tools.Workflow).GetData().(definition.Workflow)

	logrus.Debugln("终止当前流程")
	deleteCurrentStepIds := make([]int, 0)
	// 1、先终止 currentStep
	// 目标回退的步骤是 currentSteps 的前置步骤，删除 currentSteps
	currentSteps, err := this.properSet.store.FindAllCurrentStep(entry.GetEntryId())
	if err != nil {
		return err
	}
	for _, s := range currentSteps {
		if this.isPrevId(targetStep.ID, s.GetPrevIds()) {
			deleteCurrentStepIds = append(deleteCurrentStepIds, s.GetStepId())
		}
	}
	logrus.Debugf("要删除的 currentSteps: %+v", deleteCurrentStepIds)
	if len(deleteCurrentStepIds) > 0 {
		if err = this.properSet.store.DeleteCurrentSteps(deleteCurrentStepIds); err != nil {
			return err
		}
	}

	//historySteps, err := this.properSet.store.FindHistorySteps(entry.GetEntryId())
	//if err != nil {
	//	return err
	//}

	//deleteJoinTransitionPrevIds := make([]int, 0)
	//for _, h := range historySteps {
	//	if this.isPrevId(targetStep.ID, h.GetPrevIds()) {
	//		deleteJoinTransitionPrevIds = append(deleteJoinTransitionPrevIds, h.GetStepId())
	//	}
	//}
	//logrus.Debugf("要删除的 join_transitions.prev_ids: %+v", deleteJoinTransitionPrevIds)
	//if len(deleteJoinTransitionPrevIds) > 0 {
	//	if err = this.properSet.store.DeleteJoinTransitionByPrevIds(entry.GetEntryId(), deleteJoinTransitionPrevIds); err != nil {
	//		return err
	//	}
	//}

	// 2、恢复目标 step
	history, err := this.properSet.store.FindMostRecentHistory(entry.GetEntryId(), targetStep.ID)
	newStep, err := this.properSet.store.CreateCurrentStep(
		entry.GetEntryId(),
		wf,
		targetStep,
		owner,
		dueTime,
		state,
		history.GetPrevIds(),
	)
	if err != nil {
		return err
	}

	trans := NewTransition(this.properSet, this.vars, this.vars.Get(tools.CurrentStep).GetData().(spi.Step), newStep)
	return trans.Enter()

	//currentStep := this.vars.Get(tools.CurrentStep).GetData().(spi.Step)
	//logrus.Println("当前步骤：", currentStep.GetStepId(), "前置步骤:", currentStep.GetPrevIds(), "目标步骤：", targetStep.ID)
	//wf := this.vars.Get(tools.Workflow).GetData().(definition.Workflow)
	//for {
	//	step, err := this.properSet.store.FindMostRecentHistory(entry.GetEntryId(), targetStep.ID)
	//	if err != nil {
	//		break
	//	}
	//	step.GetActionName()
	//	stepDef, _ := wf.GetStep(step.GetStepId())
	//	action, _ := stepDef.GetAction(step.GetActionName())
	//	results, _ := this.getResult(action.Results)
	//	logrus.Debugf("results:%+v", results)
	//	break
	//}

	//prevIds := currentStep.GetPrevIds()
	//rollbackStepIds := make([]int, 0)
	//for i := len(prevIds) - 1; i >= 0; i-- {
	//	if prevIds[i] == targetStep.ID {
	//		break
	//	}
	//	rollbackStepIds = append(rollbackStepIds, prevIds[i])
	//}
	//
	//logrus.Println("rollbacksStepIds:", rollbackStepIds)
	return errors.New("break")
}
