package channels

import (
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// Channel is a thing on which we can send and receive messages
type Channel interface {
	Run(chan<- event.Event)
}
