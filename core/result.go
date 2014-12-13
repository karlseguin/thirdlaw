package core

type Result struct {
	Ok      bool
	Name    string
	Message string
}

func Success() *Result {
	return &Result{Ok: true, Message: ""}
}
