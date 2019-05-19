package goflow

import (
	"fmt"
	"github.com/xinzf/goflow.v2/definition"
	"github.com/xinzf/goflow.v2/tools"
)

type Workflow struct {
	transientVars *tools.TransientVars
	propertset    *PropertSet
}

func NewWorkflow(caller int, set *PropertSet) *Workflow {
	wf := &Workflow{
		propertset:    set,
		transientVars: tools.NewTransientVars(),
	}
	wf.transientVars.Put(tools.Caller, caller)
	return wf
}

func (this *Workflow) Initialize(flowId string, actionName, remark string, inputs ...map[string]interface{}) (string, error) {
	store := this.propertset.GetStore()

	// 获取工作流定义
	workflow, found := Get(flowId)
	if !found {
		return "", fmt.Errorf("Not found the workflow whith the flow_id: %s", flowId)
	}
	this.transientVars.Put(tools.Workflow, workflow)

	// 获取开始节点
	step, found := workflow.GetStartStep()
	if !found {
		return "", fmt.Errorf("Not found the workflow: %s's start step", flowId)
	}

	// 创建流程
	caller := this.transientVars.Get(tools.Caller).Int()
	entry, err := store.CreateEntry(workflow, caller)
	if err != nil {
		return "", err
	}
	this.transientVars.Put(tools.Entry, entry)

	// 创建当前步骤任务
	currentStep, err := store.CreateCurrentStep(
		entry.GetEntryId(),
		workflow,
		step,
		caller,
		definition.Result{}.GetDueTime(),
		workflow.Steps.StartInitStatus,
		[]int{},
	)
	if err != nil {
		return "", err
	}

	transition := NewTransition(this.propertset, this.transientVars, currentStep, currentStep)
	if err = transition.Enter(); err != nil {
		return "", err
	}
	if err = transition.DoAction(actionName, remark); err != nil {
		return "", err
	}

	return entry.GetEntryId(), nil
}

func (this *Workflow) DoAction(entryId string, stepId int, actionName, remark string, inputs ...map[string]interface{}) error {
	store := this.propertset.GetStore()

	entry, err := store.FindEntry(entryId)
	if err != nil {
		return err
	}
	this.transientVars.Put(tools.Entry, entry)

	//workflow, err := LoadFromString(entry.GetWorkflowXML())
	workflow, found := Get(entry.GetWorkflowId())
	if !found {
		return fmt.Errorf("Not found workflow with: %s", entry.GetWorkflowId())
	}
	//if err != nil {
	//	return err
	//}
	this.transientVars.Put(tools.Workflow, workflow)

	currentStep, found, err := store.FindCurrentStep(entryId, stepId)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("Not found current step by step_id: %d", stepId)
	}

	trans := NewTransition(this.propertset, this.transientVars, currentStep, currentStep)
	if err = trans.DoAction(actionName, remark); err != nil {
		return err
	}

	return nil
}

func (this *Workflow) GetPropertSet() *PropertSet {
	return this.propertset
}

func (this *Workflow) GetTransient() *tools.TransientVars {
	return this.transientVars
}

//func (this *Workflow) appendPrevIds(prevIds []int, currId int) []int {
//	for _, i := range prevIds {
//		if i == currId {
//			return prevIds
//		}
//	}
//
//	prevIds = append(prevIds, currId)
//	return prevIds
//}
