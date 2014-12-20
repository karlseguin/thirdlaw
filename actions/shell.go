package actions

import (
	"fmt"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"os/exec"
	"strings"
)

type Shell struct {
	command   string
	dir       string
	arguments []string
}

func (a *Shell) Run() error {
	cmd := exec.Command(a.command, a.arguments...)
	if len(a.dir) > 0 {
		cmd.Dir = a.dir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running %s %s\n   %s", a.command, strings.Join(a.arguments, " "), string(out))
	}
	return nil
}

func NewShell(t typed.Typed) *Shell {
	command, ok := t.StringIf("command")
	if ok == false {
		log.Fatalf("shell action must have a command parameter: %s", t.MustBytes(""))
	}
	return &Shell{
		command:   command,
		dir:       t.String("dir"),
		arguments: t.Strings("arguments"),
	}
}
