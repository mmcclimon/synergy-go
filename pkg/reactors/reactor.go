package reactors

import (
	"errors"
	"log"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/env"
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// Listener is a listener
type Listener = func(*event.Event)

// Reactor is a reactor...I'm still working out its interface.
type Reactor interface {
	ListenersMatching(*event.Event) []Listener
}

// Build gives you a channel based on a well-known name
func Build(name, wellKnown string, cfg config.ComponentConfig, env *env.Environment) (Reactor, error) {
	switch wellKnown {
	case "EchoReactor":
		return NewEcho(name, cfg, env), nil

	default:
		log.Fatalf("unknown reactor name %s", wellKnown)
		return nil, errors.New("unreachable")
	}
}
