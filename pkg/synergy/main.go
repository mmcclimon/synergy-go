package synergy

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Hub is the point of entry point for synergy
type Hub struct {
	name     string
	channels map[string]Channel
	reactors map[string]Reactor
	Env      *Environment
}

// NewHub gives you a new hub. Probably it will go away once I write the
// config loader.
func NewHub(name string) *Hub {
	hub := Hub{name: name}
	hub.channels = make(map[string]Channel)
	// hub.channels["slack"] = channels.NewSlack(nil)
	return &hub
}

// FromFile gives you a new hub based on a .toml file
func FromFile(filename string) *Hub {
	var config Config

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		log.Fatalf("could not read config! %s", err)
	}

	hub := Hub{
		channels: make(map[string]Channel),
		reactors: make(map[string]Reactor),
		Env:      NewEnvironment(config),
	}

	for name, cfg := range config.Channels {
		channel, _ := BuildChannel(name, cfg.Class, cfg, hub.Env)
		hub.channels[name] = channel
	}

	for name, cfg := range config.Reactors {
		reactor, _ := BuildReactor(name, cfg.Class, cfg, hub.Env)
		hub.reactors[name] = reactor
	}

	return &hub
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
