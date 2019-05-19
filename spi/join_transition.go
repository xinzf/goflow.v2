package spi

type JoinTransition interface {
	GetPrevId() int
	GetNextId() int
	GetActionName() string
	GetState() string
}
