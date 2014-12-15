package actions

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
)

var newLine = []byte("\n")

func New(t typed.Typed) core.Action {
	switch strings.ToLower(t.String("type")) {
	case "shell":
		return NewShell(t)
	default:
		panic(fmt.Errorf("invalid action type %v", string(t.MustBytes(""))))
	}
}
