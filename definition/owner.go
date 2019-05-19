package definition

import "github.com/xinzf/goflow.v2/enums"

type Owners struct {
	Owners []Owner `xml:"owner"`
}

type Owner struct {
	Type  enums.Owner `xml:"type,attr"`
	Props string      `xml:"props,attr"`
}
