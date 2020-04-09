package synergy

import (
	"fmt"
)

// EchoReactor is a reactor that echoes
type EchoReactor struct {
	name      string
	listeners map[string]Listener
}

func init() {
	RegisterReactor("EchoReactor", New)

}

// New gives you a new echo reactor
func New(name string, cfg ComponentConfig, env *Environment) Reactor {
	r := EchoReactor{
		name:      name,
		listeners: make(map[string]Listener),
	}

	r.registerHandlers()

	return &r
}

func (r *EchoReactor) registerHandlers() {
	r.registerHandler("echo", r.handleEcho, func(e *Event) bool {
		return e.WasTargeted
	})
}

func (r *EchoReactor) registerHandler(name string, handler Handler, matcher MatchFunc) {
	r.listeners[name] = Listener{
		Matcher: matcher,
		Handler: handler,
	}
}

// HandlersMatching returns a slice of listeners matching this event
func (r *EchoReactor) HandlersMatching(event *Event) []Handler {
	var handlers []Handler

	for _, listener := range r.listeners {
		if listener.Matcher(event) {
			handlers = append(handlers, listener.Handler)
		}
	}

	return handlers
}

func (r *EchoReactor) handleEcho(event *Event, errors chan<- error) {
	text := event.Text
	username := event.FromUser.Username
	event.Reply(fmt.Sprintf("I heard you, %s, when you said: %s", username, text))
}
