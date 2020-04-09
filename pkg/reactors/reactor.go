package reactors

import (
	"errors"
	"log"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/channels"
	"github.com/mmcclimon/synergy-go/pkg/env"
)

// Reactor is a reactor...I'm still working out its interface.
type Reactor interface {
	HandlersMatching(*channels.Event) []Handler
}

// Handler handles events: it's a shortcut for a unary function taking an
// event and returning void
type Handler = func(*channels.Event)

// MatchFunc takes an event and returns a bool as to whether it matches or not
type MatchFunc = func(*channels.Event) bool

// Listener is a struct that has a handler and a matchfunc
type Listener struct {
	Handler Handler
	Matcher MatchFunc
}

// ReactorBuilder is just a convenience type for a .New function
type ReactorBuilder = func(string, config.ComponentConfig, *env.Environment) Reactor

var registry = make(map[string]ReactorBuilder)

// RegisterReactor lets reactors register themselves
func RegisterReactor(wellKnown string, f ReactorBuilder) {
	log.Printf("registering reactor %s\n", wellKnown)
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
