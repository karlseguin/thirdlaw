package outputs

import (
	"bytes"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"io"
	"log"
	"strings"
	"time"
)

var newLine = []byte("\n")

type Base struct {
	next     time.Time
	output   core.Output
	snooze   time.Duration
	disabled bool
}

func (o *Base) Process(results *core.Results) {
	if o.disabled {
		return
	}
	now := time.Now()
	if now.After(o.next) {
		o.output.Process(results)
		o.next = now.Add(o.snooze)
	}
}

func New(t typed.Typed) core.Output {
	switch strings.ToLower(t.String("type")) {
	case "file":
		return build(t, NewFile(t))
	case "stdout":
		return build(t, NewStdout(t))
	case "stderr":
		return build(t, NewStderr(t))
	case "http":
		return build(t, NewHttp(t))
	default:
		log.Fatalf("invalid output type %v", string(t.MustBytes("")))
		return nil
	}
}

func build(t typed.Typed, output core.Output) core.Output {
	o := &Base{
		output:   output,
		disabled: t.BoolOr("disabled", false),
		snooze:   time.Second * time.Duration(t.IntOr("snooze", 0)),
	}
	return o
}

func writeTo(results *core.Results, writer io.Writer, newline bool) {
	writer.Write(results.Serialized())
	if newline {
		writer.Write(newLine)
	}
}

func buildBody(body []byte, r *core.Results) []byte {
	return bytes.Replace(body, []byte(`"$FRIENDLY$"`), r.Friendly(), -1)
}
