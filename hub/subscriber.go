package hub

import (
	"context"
	"errors"

	. "github.com/graph-labs-io/graphmq/types"
)

var (
	errFailedToSendMessage = errors.New("graphmq: failed to send message")
)

type Subscriber struct {
	Topic
	Recv chan (interface{})
}

// Broudcast attempts to broudcast a message to the underlying recv channel without blocking
func (s Subscriber) Broudcast(message Message) error {

	return s.broudcast(message)

}

// Broudcast attempts to broudcast a message to the underlying recv channel
func (s Subscriber) BroudcastWithTimeout(ctx context.Context, message Message) error {

	select {
	case s.Recv <- message:
		return nil
	case <-ctx.Done():
		return errFailedToSendMessage
	}
}

// Broudcast attempts to broudcast a message to the underlying recv channel
func (s Subscriber) broudcast(message Message) error {
	select {
	case s.Recv <- message:
		return nil
	default:
		return errFailedToSendMessage
	}
}
