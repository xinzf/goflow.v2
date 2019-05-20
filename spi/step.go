package spi

import (
	"github.com/xinzf/goflow.v2/definition"
	"time"
)

type Step interface {
	// 获取步骤序列号
	GetCode() string

	// 获取步骤ID
	GetStepId() int

	GetStepName() string

	// 获取实例序列号
	GetEntryId() string

	SetState(state string)

	// 获取步骤状态
	GetState() string

	// 获取步骤所有人
	GetOwner() int

	GetActionName() string

	GetActionText() string

	SetAction(action definition.Action)

	SetStartDate(t time.Time)

	// 获取步骤开始时间
	GetStartDate() time.Time

	// 获取步骤完成限制时间
	GetDueDate() time.Time

	SetFinishDate(t time.Time)

	// 获取实际完成时间
	GetFinishDate() time.Time

	SetOwner(uid int)

	SetCaller(uid int)

	// 获取实际完成人
	GetCaller() int

	SetRemark(remark string)

	GetPrevIds() []int
}
