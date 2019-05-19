package spi

type User interface {
	GetId() int
	GetName() string
	IsLeader() (bool, error)
	GetMyLeader() (User, error)
}
