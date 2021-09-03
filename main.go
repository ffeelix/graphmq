package main

import (
	"encoding/json"
	"graphmq/hub"
	. "graphmq/types"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GraphMQ struct {
	hub *hub.Hub
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	port = "80"
)

func (g *GraphMQ) HandlePublish(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var m MessageEvent

	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (g *GraphMQ) HandleSubscriber(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for {
		var s NewSubscriber
		err := ws.ReadJSON(&s)

		if err != nil {
			ws.Close()
			return
		}
		subscriber := hub.Subscriber{
			Topic: Topic(s.Topic),
		}
		g.hub.Subscribe(subscriber)

		log.Println("Sucessfully created subscriber: ", s.Topic)

	}
}

func main() {

	hub := hub.NewHub()

	mq := GraphMQ{
		hub: hub,
	}

	http.HandleFunc("/publish", mq.HandlePublish)
	http.HandleFunc("/subscribe", mq.HandleSubscriber)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
