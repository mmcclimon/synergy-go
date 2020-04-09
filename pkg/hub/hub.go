package hub

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/channels"
	"github.com/mmcclimon/synergy-go/pkg/env"
	"github.com/mmcclimon/synergy-go/pkg/reactors"

	// registry
	_ "github.com/mmcclimon/synergy-go/pkg/reactors/echo"
)

// Hub is the point of entry point for synergy
type Hub struct {
	name     string
	channels map[string]channels.Channel
	reactors map[string]reactors.Reactor
	Env      *env.Environment
}

// NewHub gives you a new hub. Probably it will go away once I write the
// config loader.
func NewHub(name string) *Hub {
	hub := Hub{name: name}
	hub.channels = make(map[string]channels.Channel)
	// hub.channels["slack"] = channels.NewSlack(nil)
	return &hub
}

// FromFile gives you a new hub based on a .toml file
func FromFile(filename string) *Hub {
	var config config.Config

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		log.Fatalf("could not read config! %s", err)
	}

	hub := Hub{
		channels: make(map[string]channels.Channel),
		reactors: make(map[string]reactors.Reactor),
		Env:      env.NewEnvironment(config),
	}

	for name, cfg := range config.Channels {
		channel, _ := channels.Build(name, cfg.Class, cfg, hub.Env)
		hub.channels[name] = channel
	}

	for name, cfg := range config.Reactors {
		reactor, _ := reactors.Build(name, cfg.Class, cfg, hub.Env)
		hub.reactors[name] = reactor
	}

	return &hub
}

// Run kicks the whole thing off. It should never exit.
func (hub *Hub) Run() {
	events := make(chan channels.Event)
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
func (hub *Hub) HandleEvent(event channels.Event, errors chan error) {
	log.Printf("%s event from %s/%s: %s",
		event.Type, event.FromChannel.Name(), event.FromUser.Username, event.Text,
	)

	var handlers []reactors.Handler
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
