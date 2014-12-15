package core

type Action interface {
	Run() error
}
