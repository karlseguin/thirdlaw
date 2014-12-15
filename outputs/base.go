package outputs

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
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
		panic(fmt.Errorf("invalid output type %v", string(t.MustBytes(""))))
	}
}

func writeTo(results *core.Results, writer io.Writer, newline bool) {
	writer.Write(results.Serialized())
	if newline {
		writer.Write(newLine)
	}
}
