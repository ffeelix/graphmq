package hub

import (
	"context"
	"time"

	. "github.com/graph-labs-io/graphmq/types"
)

type Subscriber struct {
	Topic
	Recv chan (interface{})
}

// Broudcast attempts to broudcast a message to the underlying recv channel with a timeout
func (s Subscriber) BroudcastWithTimeout(timeout time.Duration, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.broudcast(ctx, message)

}

// Broudcast attempts to broudcast a message to the underlying recv channel without blocking
func (s Subscriber) Broudcast(message interface{}) error {

	return s.broudcast(context.Background(), message)

}

// Broudcast attempts to broudcast a message to the underlying recv channel
func (s Subscriber) broudcast(ctx context.Context, message interface{}) error {
	select {
	case s.Recv <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
