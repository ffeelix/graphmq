package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	. "github.com/graph-labs-io/graphmq/types"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Message string
	Test    string
	TestStruct2
}
type TestStruct2 struct {
	Message string
	Test    string
}

func TestPublish(t *testing.T) {

	message := MessageEvent{
		Topic: "Test",
		Message: TestStruct{
			"",
			"",
			TestStruct2{"", ""},
		},
		Method:      FanoutMessage,
		IsPersisted: true,
	}
	b, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", "http://localhost:80/publish", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	assert.Equal(t, 200, resp.StatusCode)

}
