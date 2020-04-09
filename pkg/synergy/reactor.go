package synergy

import (
	"errors"
	"log"
)

// GenericReactor is a struct suitable for embedding in actual reactors. Doing
// so gives you some generic behavior
type GenericReactor struct {
	Name      string
	Listeners map[string]ReactorListener
}

// Reactor is a reactor...I'm still working out its interface.
type Reactor interface {
	HandlersMatching(*Event) []Handler
}

// Handler handles events: it's a shortcut for a unary function taking an
// event and returning void
type Handler = func(*Event, chan<- error)

// MatchFunc takes an event and returns a bool as to whether it matches or not
type MatchFunc = func(*Event) bool

// ReactorBuilder is just a convenience type for a builder function
type ReactorBuilder = func(string, ComponentConfig, *Environment) Reactor

// ReactorListener is a struct that has a handler and a matchfunc
type ReactorListener struct {
	Handler Handler
	Matcher MatchFunc
}

// used by BuildReactor below
var reactorRegistry = make(map[string]ReactorBuilder)

// RegisterReactor lets reactors register themselves
func RegisterReactor(wellKnown string, f ReactorBuilder) {
	log.Printf("registering reactor %s\n", wellKnown)
	reactorRegistry[wellKnown] = f
}

// BuildReactor gives you a reactor based on a well-known name
func BuildReactor(name, wellKnown string, cfg ComponentConfig, env *Environment) (Reactor, error) {
	builder, ok := reactorRegistry[wellKnown]

	if !ok {
		log.Fatalf("unknown reactor name %s", wellKnown)
		return nil, errors.New("unreachable")
	}

	return builder(name, cfg, env), nil
}

func (r *GenericReactor) registerHandler(name string, handler Handler, matcher MatchFunc) {
	r.Listeners[name] = ReactorListener{
		Matcher: matcher,
		Handler: handler,
	}
}

// HandlersMatching returns a slice of listeners matching this event
func (r *GenericReactor) HandlersMatching(event *Event) []Handler {
	var handlers []Handler

	for _, listener := range r.Listeners {
		if listener.Matcher(event) {
			handlers = append(handlers, listener.Handler)
		}
	}

	return handlers
}
