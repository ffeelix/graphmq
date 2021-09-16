package hub

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	. "github.com/graph-labs-io/graphmq/types"
)

// Hub handles all subscriber & producer interactions
type Hub struct {

	// Maps each topic to a list of subscribers
	subscribers map[Subscriber]bool

	// register is the underlying channel to signify to in order to create a new subscriber
	register chan (Subscriber)

	// unregister is the underlying channel to signify to in order to remove a subscriber
	unregister chan (Subscriber)

	broudcast chan (MessageEvent)

	deleteCache chan (Topic)

	cache map[Topic]Message
}

// NewHub starts a new hub while also starting a goroutine
func NewHub() *Hub {

	h := &Hub{
		subscribers: make(map[Subscriber]bool),

		register:   make(chan Subscriber),
		unregister: make(chan Subscriber),
		broudcast:  make(chan MessageEvent),
		cache:      make(map[Topic]Message),
	}
	go h.worker()
	return h

}

// fanout sends the given message to all appropiate subscribers that match the message's topic
func (h *Hub) fanout(message MessageEvent) {
	for subscriber := range h.subscribers {
		if strings.EqualFold(string(subscriber.Topic), string(message.Topic)) {
			if err := subscriber.Broudcast(message.Message); err != nil {
				log.Println(err.Error())
			}
		}
	}

}

// singleBroudcast sends the given message to the first subscribers that match the message's topic
func (h *Hub) singleBroudcast(message MessageEvent) {
	for subscriber := range h.subscribers {
		if strings.EqualFold(string(subscriber.Topic), string(message.Topic)) {

			if err := subscriber.Broudcast(message.Message); err != nil {
				log.Println(err.Error())
				continue

			}
			return
		}
	}

}

// runs in a seperate goroutine
func (h *Hub) worker() {
	for {
		select {
		case subscriber := <-h.register:

			// Create a new entry with the specified subscriber

			// We can just send back the cached answer if it exists
			if message, ok := h.cache[subscriber.Topic]; ok {
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
					defer cancel()

					if err := subscriber.BroudcastWithTimeout(ctx, message); err != nil {
						log.Println(err.Error())
					}
				}()

				continue
				// log.Println("Successfully message pulled from cache")
			}

			h.subscribers[subscriber] = true
			log.Println(fmt.Sprintf("Sucessfully created Subscriber: %s | Total Subs: %d", subscriber.Topic, len(h.subscribers)))

		case subscriber := <-h.unregister:

			delete(h.subscribers, subscriber)
			log.Println(fmt.Sprintf("Sucessfully Deleted Subscriber: %s | Total Subs: %d", subscriber.Topic, len(h.subscribers)))

		case messageEvent := <-h.broudcast:
			switch messageEvent.Method {
			case FanoutMessage:
				h.fanout(messageEvent)
			case SingleReceiver:
				h.singleBroudcast(messageEvent)
			}

			if messageEvent.IsPersisted {
				// If the message requires persistence then we can add it to the local cache
				h.cache[messageEvent.Topic] = messageEvent.Message

			}
		}
	}
}

func (h *Hub) DeleteCache(t Topic) {
	h.deleteCache <- t

}

// Unsubscribe removes the specified subscriber from the hub
func (h *Hub) Unsubscribe(subscriber Subscriber) {
	h.unregister <- subscriber
}

// Subscribe adds a the specified subscriber into the hub
func (h *Hub) Subscribe(subscriber Subscriber) {
	h.register <- subscriber
}

// Broudcast sends the given message into the broudcast channel
func (h *Hub) Broudcast(message MessageEvent) {

	h.broudcast <- message

}
