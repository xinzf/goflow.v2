package enums

type EntryState int

const (
	UNKNOWN EntryState = iota - 1
	CREATED
	ACTIVATED
	SUSPENDED
	KILLED
	COMPLETED
)

func (static EntryState) Int() int {
	return int(static)
}

func (static EntryState) Text() string {
	m := map[EntryState]string{
		UNKNOWN:   "未知状态",
		CREATED:   "已创建",
		ACTIVATED: "进行中",
		SUSPENDED: "暂停",
		KILLED:    "终止",
		COMPLETED: "完成",
	}

	s, _ := m[static]
	return s
}
