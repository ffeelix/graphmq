package main_test

import (
	"log"
	"testing"

	. "github.com/graph-labs-io/graphmq/types"

	"github.com/gorilla/websocket"
)

func TestSync(t *testing.T) {

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:80/subscribe", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = c.WriteJSON(NewSubscriber{Topic: "Test"})
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		_, m, _ := c.ReadMessage()
		log.Println(string(m))
	}
}
