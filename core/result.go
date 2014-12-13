package core

import (
	"encoding/json"
	"log"
)
var serializedError = []byte(`{"error": "failed to serialize results"}`)

type Result struct {
	Ok      bool
	Name    string
	Message string
}

func Success() *Result {
	return &Result{Ok: true, Message: ""}
}

type Results struct {
	serialized []byte
	List []*Result
}

func NewResults(list []*Result) *Results {
	return &Results{List: list}
}

func (r *Results) Serialized() []byte {
	if r.serialized == nil {
		data, err := json.Marshal(r.List)
		if err != nil {
			log.Println("failed to serialize results", err)
			data = serializedError
		}
		r.serialized = data
	}
	return r.serialized
}
