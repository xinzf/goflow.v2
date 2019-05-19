package enums

type Owner string

const (
	//AllUsers    OperatorType = "all"
	Creator  Owner = "creator"
	Caller   Owner = "caller"
	Leader   Owner = "leader"
	Users    Owner = "users"
	Deps     Owner = "deps"
	Roles    Owner = "roles"
	Variable Owner = "variable"
	None     Owner = "none"
)

func (static Owner) String() string {
	return string(static)
}

func (static Owner) Text() string {
	m := map[Owner]string{
		//AllUsers:    "所有人",
		Creator:  "流程发起人",
		Caller:   "当前节点负责人",
		Leader:   "上级领导",
		Users:    "指定用户",
		Deps:     "指定部门",
		Roles:    "指定角色",
		Variable: "指定变量",
	}

	s, _ := m[static]
	return s
}
