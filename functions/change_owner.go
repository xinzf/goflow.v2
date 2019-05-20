package functions

import (
	"errors"
	"github.com/xinzf/goflow.v2/spi"
	"github.com/xinzf/goflow.v2/tools"
	"log"
)

type ChangeOwner struct {
}

func (this *ChangeOwner) GetName() string {
	return "changeOwner"
}

func (this *ChangeOwner) Eval(store spi.Store, transientVars *tools.TransientVars, args map[string]*tools.Value) (interface{}, error) {
	log.Printf("%+v",args)
	val, ok := args["uid"]
	if !ok {
		return nil, errors.New("function: changeOwner 缺少参数: task_uid")
	}

	var uid int
	if key, flag := val.ParseVariable(); flag {
		uid = transientVars.Get(key).Int()
	} else {
		uid = val.Int()
	}

	if uid == 0 {
		return nil, errors.New("function: changeOwner 参数 task_uid 不能为0")
	}
	currentStep := transientVars.Get(tools.CurrentStep).GetData().(spi.Step)
	currentStep.SetOwner(uid)

	if err := store.UpdateCurrentStep(currentStep); err != nil {
		return nil, err
	}
	return nil, nil
}
