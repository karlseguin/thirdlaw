package core

var OKResult = &Result{true, ""}

type Result struct {
	Ok      bool
	Message string
}
