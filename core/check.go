package core

type Runner interface {
	Run() *Result
}

type Check interface {
	Runner
	Name() string
}
