package definition

import "time"

type Results struct {
	Default []Result `xml:"default-result"`
	Results []Result `xml:"result"`
}

type Result struct {
	Step       string     `xml:"step,attr"`
	Status     string     `xml:"status,attr"`
	ExitStatus string     `xml:"exit-status,attr"`
	DueSeconds int64      `xml:"due-seconds,attr"`
	Split      int        `xml:"split,attr"`
	Join       int        `xml:"join,attr"`
	Owners     Owners     `xml:"owners"`
	Conditions Conditions `xml:"conditions"`
}

func (static Result) GetDueTime() time.Time {
	due := int64(31536000)
	if static.DueSeconds > 0 {
		due = static.DueSeconds
	}
	return time.Now().Add(time.Duration(due) * time.Second)
}
