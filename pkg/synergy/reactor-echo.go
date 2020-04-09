package synergy

import (
	"fmt"
)

// EchoReactor is a reactor that echoes
type EchoReactor struct {
	GenericReactor
}

func init() {
	RegisterReactor("EchoReactor",
		func(name string, cfg ComponentConfig, env *Environment) Reactor {
			r := EchoReactor{
				GenericReactor{
					Name:      name,
					Listeners: make(map[string]ReactorListener),
				},
			}

			r.registerHandlers()

			return &r
		},
	)
}

func (r *EchoReactor) registerHandlers() {
	r.registerHandler("echo", r.handleEcho, func(e *Event) bool {
		return e.WasTargeted
	})
}

func (r *EchoReactor) handleEcho(event *Event, errors chan<- error) {
	text := event.Text
	username := event.FromUser.Username
	event.Reply(fmt.Sprintf("I heard you, %s, when you said: %s", username, text))
}
