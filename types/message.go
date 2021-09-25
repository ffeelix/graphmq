package types

import "time"

var (
	// FanoutMessage fans out the recieved message to all associated subscribers
	FanoutMessage = Method("FANOUT_MESSAGE")

	// SingleReceiver only sends the message to the first subscriber found
	SingleReceiver = Method("RECEIVER_MESSAGE")
)

type Method string
type Message interface{}

type MessageEvent struct {
	Topic       Topic   `json:"topic"`
	Message     Message `json:"message"`
	Method      Method  `json:"method"`
	IsPersisted bool    `json:"isPersisted"`

	ExpiryDate time.Time `json:"expiryTime,omitempty"`
}

func (m MessageEvent) IsExpired() bool {
	return time.Now().After(m.ExpiryDate)
}
