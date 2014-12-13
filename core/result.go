package core

import (
	"encoding/json"
	"log"
	"time"
)

var serializedError = []byte(`{"error": "failed to serialize results"}`)

type Result struct {
	Ok       bool   `json:"ok"`
	Name     string `json:"name"`
	Message  string `json:"message"`
	Milliseconds int    `json:"ms"`
}

func Success() *Result {
	return &Result{Ok: true, Message: ""}
}

type Results struct {
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
