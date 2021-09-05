package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/graph-labs-io/graphmq/hub"
	. "github.com/graph-labs-io/graphmq/types"

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

	port             = "80"
	healthCheckTimer = 5 * time.Second
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

	g.hub.Broudcast(m)

}

// HandleSubscriber creates a new subscriber
func (g *GraphMQ) HandleSubscriber(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer ws.Close()

	for {
		var s NewSubscriber
		err := ws.ReadJSON(&s)

		if err != nil {

			return
		}
		log.Println("Sucessfully created Subscriber:", s.Topic)

		recv := make(chan (interface{}))

		subscriber := hub.Subscriber{
			Topic: s.Topic,
			Recv:  recv,
		}
		g.hub.Subscribe(subscriber)
		defer g.hub.Unsubscribe(subscriber)

		timer := time.NewTimer(healthCheckTimer)
		for {
			timer.Reset(healthCheckTimer)
			select {
			case m := <-recv:
				ws.WriteJSON(m)
			case <-timer.C:
				err := ws.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Println("Failed health check")
					return
				}
			}
		}

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
