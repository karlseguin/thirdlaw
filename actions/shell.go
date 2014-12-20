package actions

import (
	"fmt"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"os"
	"strings"
	"syscall"
)

type Shell struct {
	command   string
	dir       string
	arguments []string
}

func (a *Shell) Run() error {
	_, err := os.StartProcess(a.command, a.arguments, &os.ProcAttr{Dir: a.dir, Sys: &syscall.SysProcAttr{Setpgid: true}})
	if err != nil {
		return fmt.Errorf("error running %s %s", a.command, strings.Join(a.arguments, " "))
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
