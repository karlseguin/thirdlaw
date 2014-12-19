package core

type ActionProvider interface {
	GetAction(name string) Action
}

type Action interface {
	Run() error
}
