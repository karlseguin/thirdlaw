package actions

import (
	"bytes"
	"fmt"
	"gopkg.in/karlseguin/typed.v1"
	"os/exec"
	"strings"
)

type Shell struct {
	command   string
	arguments []string
}

func (a *Shell) Run() error {
	var out bytes.Buffer
	cmd := exec.Command(a.command, a.arguments...)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s %s\n   %s", a.command, strings.Join(a.arguments, " "), out.String())
	}
	return nil
}

func NewShell(t typed.Typed) *Shell {
	command, ok := t.StringIf("command")
	if ok == false {
		panic(fmt.Sprintf("shell action must have a command parameter: %s", t.MustBytes("")))
	}
	return &Shell{
		command:   command,
		arguments: t.Strings("arguments"),
	}
}
