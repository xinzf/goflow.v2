package enums

type State string

const (
	Queue      State = "queue"
	Doing      State = "doing"
	Done       State = "done"
	Draft      State = "draft"
	Supplement State = "supplement"
)

func (static State) String() string {
	return string(static)
}

func (static State) Text() string {
	m := map[State]string{
		Queue:      "待处理",
		Doing:      "处理中",
		Done:       "完毕",
		Draft:      "录入中",
		Supplement: "待完善",
	}

	s, _ := m[static]
	return s
}
