package enums

type ConditionRelation string

const (
	AND ConditionRelation = "AND"
	OR  ConditionRelation = "OR"
)

func (static ConditionRelation) String() string {
	return string(static)
}
