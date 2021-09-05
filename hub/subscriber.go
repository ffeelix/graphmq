package hub

import (
	"context"

	. "github.com/graph-labs-io/graphmq/types"
)

type Subscriber struct {
	Topic
	Recv chan (interface{})
}

// Broudcast attempts to broudcast a message to the underlying recv channel without blocking
func (s Subscriber) Broudcast(message Message) error {

	return s.broudcast(context.Background(), message)

}

// Broudcast attempts to broudcast a message to the underlying recv channel
func (s Subscriber) broudcast(ctx context.Context, message Message) error {
	select {
	case s.Recv <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
