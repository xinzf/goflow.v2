package spi

import (
	"github.com/xinzf/goflow.v2/enums"
	"time"
)

type Entry interface {
	// 获取实例ID
	GetEntryId() string

	SetState(state enums.EntryState)

	GetState() enums.EntryState

	SetStartTime(t time.Time)

	GetStartTime() time.Time

	SetEndTime(t time.Time)

	GetEndTime() time.Time

	GetCreator() int

	// 获取工作流序列号
	GetWorkflowName() string

	GetWorkflowId() string

	GetDescribe() string

	GetWorkflowXML() string
}
