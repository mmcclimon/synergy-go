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
	name      string
	listeners map[string]reactors.Listener
}

func init() {
	reactors.RegisterReactor("EchoReactor", New)

}

// New gives you a new echo reactor
func New(name string, cfg config.ComponentConfig, env *env.Environment) reactors.Reactor {
	r := EchoReactor{
		name:      name,
		listeners: make(map[string]reactors.Listener),
	}

	r.registerHandlers()

	return &r
}

func (r *EchoReactor) registerHandlers() {
	r.registerHandler("echo", r.handleEcho, func(e *channels.Event) bool {
		return e.WasTargeted
	})
}

func (r *EchoReactor) registerHandler(name string, handler reactors.Handler, matcher reactors.MatchFunc) {
	r.listeners[name] = reactors.Listener{
		Matcher: matcher,
		Handler: handler,
	}
}

// HandlersMatching returns a slice of listeners matching this event
func (r *EchoReactor) HandlersMatching(event *channels.Event) []reactors.Handler {
	var handlers []reactors.Handler

	for _, listener := range r.listeners {
		if listener.Matcher(event) {
			handlers = append(handlers, listener.Handler)
		}
	}

	return handlers
}

func (r *EchoReactor) handleEcho(event *channels.Event) {
	text := event.Text
	username := event.FromUser.Username
	event.Reply(fmt.Sprintf("I heard you, %s, when you said: %s", username, text))
}
