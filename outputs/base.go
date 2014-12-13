package outputs

import (
	"fmt"
	"github.com/karlseguin/beats/core"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
)

func New(t typed.Typed) core.Output {
	switch strings.ToLower(t.String("type")) {
	case "file":
		return NewFile(t)
	default:
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("invalid output type %v", string(b)))
	}
}
