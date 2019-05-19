package definition

import "github.com/xinzf/goflow.v2/enums"

type Conditions struct {
	Type      enums.ConditionRelation `xml:"type,attr"`
	Functions []Function              `xml:"function"`
}
