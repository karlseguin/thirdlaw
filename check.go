package beats

import (
	"fmt"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
)

var OKResult = &Result{true, ""}

type Named interface{
	Name() string
}

type Check interface {
	Named
	Run() *Result
}

func CheckFactory(t typed.Typed) Check {
	name, ok := t.StringIf("name")
	if ok == false {
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("missing name %v", string(b)))
	}
	switch strings.ToLower(t.String("type")) {
	case "http":
		return NewHttpCheck(name, t)
	default:
		panic(fmt.Errorf("unknown type for %q", name))
	}
}
