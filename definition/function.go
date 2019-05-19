package definition

import "github.com/xinzf/goflow.v2/tools"

type Function struct {
	Name string `xml:"name,attr"`
	Args []Arg  `xml:"arg"`
	Want string `xml:"want,attr"`
}

func (static Function) GetWant() bool {
	if static.Want == "false" || static.Want == "False" || static.Want == "FALSE" {
		return false
	} else {
		return true
	}
}

func (static Function) GetArgValues() map[string]*tools.Value {
	vals := make(map[string]*tools.Value)
	for _, a := range static.Args {
		vals[a.Name] = tools.NewValue(a.Value)
	}
	return vals
}
