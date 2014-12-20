package outputs

import (
	"bytes"
	"encoding/json"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Http struct {
	address string
	timeout time.Duration
	body    interface{}
	client  *http.Client
}

func (o Http) Process(results *core.Results) {
	method := "GET"
	var body io.Reader
	if o.body != nil {
		method = "POST"
		body = bytes.NewBuffer(buildBody(o.body.([]byte), results))
	}
	req, err := http.NewRequest(method, o.address, body)
	if err != nil {
		log.Println("failed to create http request", err)
		return
	}
	resp, err := o.client.Do(req)
	if err != nil {
		log.Println("failed http output", o.address, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println("http output non 2xx response", string(data))
	}
}

func NewHttp(t typed.Typed) *Http {
	address, ok := t.StringIf("address")
	if ok == false {
		log.Fatal("http output requires an address")
	}
	body := t["body"]
	if body != nil {
		var err error
		body, err = json.Marshal(body)
		if err != nil {
			log.Fatal("failed to prepare output body", err)
		}
	}

	return &Http{
		address: address,
		body:    body,
		client: &http.Client{
			Timeout: time.Second * time.Duration(t.IntOr("timeout", 5)),
		},
	}
}
