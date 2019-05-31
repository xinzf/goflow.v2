package tools

import (
	"github.com/emirpasic/gods/maps/hashmap"
)

const (
	Entry       string = "entry"
	PrevStep    string = "prev_step"
	CurrentStep string = "current_step"
	Workflow    string = "workflow"
	Caller      string = "caller"
	Inputs      string = "inputs"
	NewSteps    string = "new_steps"
)

func NewTransientVars() *TransientVars {
	return &TransientVars{data: hashmap.New()}
}

type TransientVars struct {
	data *hashmap.Map
}

func (this *TransientVars) Put(key string, val interface{}) {
	this.data.Put(key, NewValue(val))
}

func (this *TransientVars) Get(key string) *Value {
	//logrus.Println(key)
	if v, found := this.data.Get(key); found {
		return v.(*Value)
	}
	return &Value{}
}
