package reactors

import (
	"fmt"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/env"
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// EchoReactor is a reactor that echoes
type EchoReactor struct {
	name string
}

// NewEcho gives you a new echo reactor
func NewEcho(name string, cfg config.ComponentConfig, env *env.Environment) *EchoReactor {
	self := EchoReactor{
		name: name,
	}

	return &self
}

func (r *EchoReactor) ListenersMatching(*event.Event) []Listener {
	handlers := make([]func(*event.Event), 0)
	handlers = append(handlers, r.handleEcho)

	return handlers
}

func (r *EchoReactor) handleEcho(*event.Event) {
	fmt.Println("would echo")
}
