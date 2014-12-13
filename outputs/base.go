package outputs

import (
	"fmt"
	"github.com/karlseguin/beats/core"
	"gopkg.in/karlseguin/typed.v1"
	"io"
	"strings"
)

var newLine = []byte("\n")

func New(t typed.Typed) core.Output {
	switch strings.ToLower(t.String("type")) {
	case "file":
		return NewFile(t)
	case "stdout":
		return NewStdout(t)
	case "stderr":
		return NewStderr(t)
	default:
		b, _ := t.ToBytes("")
		panic(fmt.Errorf("invalid output type %v", string(b)))
	}
}

func writeTo(results *core.Results, writer io.Writer, newline bool) {
	writer.Write(results.Serialized())
	if newline {
		writer.Write(newLine)
	}
}
