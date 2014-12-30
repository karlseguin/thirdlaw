package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var serializedError = []byte(`{"error": "failed to serialize results"}`)

type Result struct {
	Ok           bool   `json:"ok"`
	Silent       bool   `json:"silent"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Milliseconds int    `json:"ms"`
	Failures     int    `json:"failures"`
}

func Success() *Result {
	return &Result{Ok: true, Message: ""}
}

func Error(err error) *Result {
	return &Result{Ok: false, Message: err.Error()}
}

func Failuref(format string, args ...interface{}) *Result {
	return &Result{Ok: false, Message: fmt.Sprintf(format, args...)}
}

type Results struct {
	friendly   []byte
	serialized []byte
	Time       time.Time `json:"time"`
	List       []*Result `json:"results"`
}

func NewResults(list []*Result) *Results {
	return &Results{
		List: list,
		Time: time.Now().UTC(),
	}
}

func (r *Results) Serialized() []byte {
	if r.serialized == nil {
		data, err := json.Marshal(r)
		if err != nil {
			log.Println("failed to serialize results", err)
			data = serializedError
		}
		r.serialized = data
	}
	return r.serialized
}

func (r *Results) Friendly() []byte {
	if len(r.friendly) == 0 {
		data := bytes.NewBuffer(make([]byte, 0, 2048))
		for _, res := range r.List {
			data.WriteString("* " + res.Name + "\t")
			if res.Ok {
				data.WriteString("OK\n")
			} else {
				data.WriteString("FAIL\n\t" + res.Message + "\n")
			}
		}
		r.friendly, _ = json.Marshal(data.String())
	}
	return r.friendly
}
