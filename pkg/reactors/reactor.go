package reactors

import (
	"errors"
	"fmt"
	"log"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/channels"
	"github.com/mmcclimon/synergy-go/pkg/env"
)

// Listener is a listener
type Listener = func(*channels.Event)

// Reactor is a reactor...I'm still working out its interface.
type Reactor interface {
	ListenersMatching(*channels.Event) []Listener
}

// ReactorBuilder is just a convenience type for a .New function
type ReactorBuilder = func(string, config.ComponentConfig, *env.Environment) Reactor

var registry = make(map[string]ReactorBuilder)

// RegisterReactor lets reactors register themselves
func RegisterReactor(wellKnown string, f ReactorBuilder) {
	fmt.Printf("registering reactor %s\n", wellKnown)
	registry[wellKnown] = f
}

// Build gives you a channel based on a well-known name
func Build(name, wellKnown string, cfg config.ComponentConfig, env *env.Environment) (Reactor, error) {
	builder, ok := registry[wellKnown]

	if !ok {
		log.Fatalf("unknown reactor name %s", wellKnown)
		return nil, errors.New("unreachable")
	}

	return builder(name, cfg, env), nil
}
