package checks

import (
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"fmt"
	"os/exec"
)

type Shell struct {
	command   string
	dir       string
	arguments []string
	out       interface{}
}

func (c *Shell) Run() *core.Result {
	cmd := exec.Command(c.command, c.arguments...)
	if len(c.dir) > 0 {
		cmd.Dir = c.dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return core.Error(err)
	}
	if c.out != nil && c.out.(string) != string(out) {
		return core.Failuref("Expected output %q got %q", c.out.(string), string(out))
	}
	return core.Success()
}

func NewShell(t typed.Typed) *Shell {
	command, ok := t.StringIf("command")
	if ok == false {
		panic(fmt.Sprintf("shell check must have a command parameter: %s", t.MustBytes("")))
	}
	return &Shell{
		command:   command,
		dir:       t.String("dir"),
		arguments: t.Strings("arguments"),
		out:    t["out"],
	}
}
