package echo

import (
	"fmt"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/channels"
	"github.com/mmcclimon/synergy-go/pkg/env"
	"github.com/mmcclimon/synergy-go/pkg/reactors"
)

// EchoReactor is a reactor that echoes
type EchoReactor struct {
	name string
}

func init() {
	reactors.RegisterReactor("EchoReactor", New)
}

// New gives you a new echo reactor
func New(name string, cfg config.ComponentConfig, env *env.Environment) reactors.Reactor {
	self := EchoReactor{
		name: name,
	}

	return &self
}

// ListenersMatching returns a slice of listeners matching this event
func (r *EchoReactor) ListenersMatching(event *channels.Event) []reactors.Listener {
	handlers := make([]func(*channels.Event), 0)
	handlers = append(handlers, r.handleEcho)

	return handlers
}

func (r *EchoReactor) handleEcho(event *channels.Event) {
	text := event.Text
	username := event.FromUser.Username
	event.Reply(fmt.Sprintf("I heard you, %s, when you said: %s", username, text))
}
