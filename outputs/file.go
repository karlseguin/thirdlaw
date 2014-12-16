package outputs

import (
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"os"
)

type File struct {
	path     string
	flags    int
}

func (o File) Process(results *core.Results) {
	f, err := os.OpenFile(o.path, o.flags, 0600)
	if err != nil {
		log.Println("failed to open output file", o.path, err)
		return
	}
	defer f.Close()
	writeTo(results, f, true)
}

func NewFile(t typed.Typed) *File {
	flags := os.O_CREATE | os.O_WRONLY
	if t := t.BoolOr("truncate", false); t {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	return &File{
		flags: flags,
		path: t.StringOr("path", "failures.log"),
	}
}
