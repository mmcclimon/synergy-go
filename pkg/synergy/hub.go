package synergy

import "log"

// Hub is the point of entry for synergy
type Hub struct {
	name     string
	channels map[string]Channel
	reactors map[string]Reactor
	Env      *Environment
}

// Run kicks the whole thing off. It should never exit.
func (hub *Hub) Run() {
	events := make(chan Event)
	errors := make(chan error)

	for name, channel := range hub.channels {
		log.Printf("starting channel %s\n", name)
		go channel.Run(events)
	}

	for event := range events {
		hub.HandleEvent(event, errors)
	}
}

// HandleEvent handles events, yo
func (hub *Hub) HandleEvent(event Event, errors chan error) {
	log.Printf("%s event from %s/%s: %s",
		event.Type, event.FromChannel.Name(), event.FromUser.Username, event.Text,
	)

	var handlers []Handler
	for _, reactor := range hub.reactors {
		handlers = append(handlers, reactor.HandlersMatching(&event)...)
	}

	for _, handler := range handlers {
		go handler(&event, errors)
	}

	for err := range errors {
		log.Println("caught error with reactor:", err)
	}
}
