package spi

import (
	"github.com/xinzf/goflow.v2/definition"
	"time"
)

type Store interface {
	// 查询实例
	FindEntry(entryId string) (Entry, error)

	// 创建实例
	CreateEntry(workflow definition.Workflow, owner int) (Entry, error)

	UpdateEntry(entry Entry) error

	// 查询当前步骤
	FindCurrentStep(entryId string, stepId int) (Step, bool, error)

	FindAllCurrentStep(entryId string) ([]Step, error)

	// 创建当前步骤
	CreateCurrentStep(entryId string, workflow definition.Workflow, step definition.Step, owner int, dueDate time.Time, state string, prevIds []int) (Step, error)

	UpdateCurrentStep(step Step) error

	DeleteAllCurrentStep(entryId string) error

	DeleteCurrentSteps(stepIds []int) error

	CreateCopy(entryId string, uids []int) error

	// 查询历史步骤
	FindHistorySteps(entryId string) ([]Step, error)

	FindMostRecentHistory(entryId string, stepId int) (Step, error)

	// 移动步骤到历史记录
	MoveHistory(step Step) error

	GetUser(uid int) (User, error)

	GetUsers(uids []int) ([]User, error)

	GetUsersByDepIds(depIds []int) ([]User, error)

	GetUsersByRoleIds(roleIds []int) ([]User, error)

	CreateJoinTransition(entryId string, currentStep Step, nextId int) error

	GetJoinTransitionsByNextId(entryId string, nextId int) ([]JoinTransition, error)

	GetJoinTransitionsByPrevId(entryId string, prevId int) (JoinTransition, error)

	DeleteJoinTransition(entryId string, nextStepId int) error

	DeleteJoinTransitionByPrevIds(entryId string, prevIds []int) (err error)
}
