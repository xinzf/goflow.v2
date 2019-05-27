package definition

import (
	"encoding/json"
	"encoding/xml"
	"github.com/pkg/errors"
	"reflect"
)

type Workflow struct {
	ID       string  `xml:"id,attr"`
	Version  string  `xml:"version,attr"`
	Describe string  `xml:"describe,attr"`
	Name     string  `xml:"name,attr"`
	Prefix   string  `xml:"prefix,attr"`
	Steps    Steps   `xml:"steps"`
	Joins    Joins   `xml:"joins"`
	Splits   Splits  `xml:"splits"`
	Extends  Extends `xml:"extends"`
}

func (static Workflow) GetStartStep() (Step, bool) {
	return static.GetStep(static.Steps.Start)
}

func (static Workflow) BindExtends(obj interface{}) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("Extend 不是有效的指针")
	}

	//log.Println(static.Extends)
	return xml.Unmarshal([]byte(static.Extends.Text), obj)
	return nil
}

func (static Workflow) GetStep(id int) (Step, bool) {
	for _, s := range static.Steps.Steps {
		if s.ID == id {
			return s, true
		}
	}
	return Step{}, false
}

func (static Workflow) GetSteps() []Step {
	return static.Steps.Steps
}

func (static Workflow) GetStepsCount() int {
	return len(static.GetSteps())
}

func (static Workflow) GetJoins() []Join {
	return static.Joins.Join
}

func (static Workflow) GetJoin(id int) (Join, bool) {
	for _, j := range static.Joins.Join {
		if j.Id == id {
			return j, true
		}
	}
	return Join{}, false
}

func (static Workflow) GetSplits() []Split {
	return static.Splits.Split
}

func (static Workflow) GetSplit(id int) (Split, bool) {
	for _, s := range static.Splits.Split {
		if s.Id == id {
			return s, true
		}
	}
	return Split{}, false
}

func (static Workflow) ToXML() string {
	b, _ := xml.Marshal(static)
	return string(b[:])
}

func (static Workflow) ToJson() string {
	b, _ := json.Marshal(static)
	return string(b[:])
}
