package definition

type Steps struct {
	Start           int    `xml:"start,attr"`
	End             int    `xml:"end,attr"`
	Steps           []Step `xml:"step"`
	StartInitStatus string `xml:"start-init-status,attr"`
}

type Step struct {
	ID           int     `xml:"id,attr"`
	Name         string  `xml:"name,attr"`
	Actions      Actions `xml:"actions"`
	PreFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"pre-functions"`
	PostFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"post-functions"`
	Restrict Restrict `xml:"restrict"`
}

func (static Step) ActionNames() []string {
	names := make([]string, 0)
	for _, v := range static.Actions.Actions {
		names = append(names, v.Name)
	}

	return names
}

func (static Step) GetAction(name string) (Action, bool) {
	for _, v := range static.Actions.Actions {
		if v.Name == name {
			return v, true
		}
	}
	return Action{}, false
}

func (static Step) GetAutoActions() []Action {

	actions := make([]Action, 0)
	for _, a := range static.Actions.Actions {
		if a.Auto == true {
			actions = append(actions, a)
		}
	}
	return actions

}

//func (static Step) GetDueTime(startTime time.Time) time.Time {
//	due := int64(31536000)
//	if static.DueTime > 0 {
//		due = static.DueTime
//	}
//	return startTime.Add(time.Duration(due) * time.Second)
//}
