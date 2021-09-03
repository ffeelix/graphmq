package types

var (
	// FanoutMessage fans out the recieved message to all associated subscribers
	FanoutMessage = Method("FANOUT_MESSAGE")

	// SingleReceiver only sends the message to the first subscriber found
	SingleReceiver = Method("RECEIVER_MESSAGE")
)

type Method string

type MessageEvent struct {
	Topic       Topic       `json:"topic"`
	Message     interface{} `json:"message"`
	Method      Method      `json:"method"`
	IsPersisted bool        `json:"isPersisted"`
}
