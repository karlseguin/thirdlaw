package outputs

import (
	"encoding/json"
	"github.com/karlseguin/beats/core"
	"gopkg.in/karlseguin/typed.v1"
	"log"
	"os"
)

type File struct {
	path     string
	truncate bool
}

func (o *File) Process(results []*core.Result) {
	data, err := json.Marshal(results)
	if err != nil {
		log.Println("file output failed to serialize resuls", err)
		return
	}
	f, err := os.OpenFile(o.path, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0600)
	if err != nil {
		log.Println("failed to open output file", o.path, err)
		return
	}
	defer f.Close()
	f.Write(data)
}

func NewFile(t typed.Typed) *File {
	return &File{
		path:     t.StringOr("path", "failures.log"),
	}
}
