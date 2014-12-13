package outputs

import (
	"github.com/karlseguin/beats/core"
	"gopkg.in/karlseguin/typed.v1"
	"os"
)

type Stdout struct{}

func (_ Stdout) Process(results *core.Results) {
	writeTo(results, os.Stdout, true)
}

func NewStdout(t typed.Typed) *Stdout {
	return &Stdout{}
}

type Stderr struct{}

func (_ Stderr) Process(results *core.Results) {
	writeTo(results, os.Stderr, true)
}

func NewStderr(t typed.Typed) *Stderr {
	return &Stderr{}
}
