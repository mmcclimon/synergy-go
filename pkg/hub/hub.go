package hub

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/channels"
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// Hub is the point of entry point for synergy
type Hub struct {
	name     string
	channels map[string]channels.Channel
	Env      *Environment
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

	md, err := toml.DecodeFile(filename, &config)
	if err != nil {
		log.Fatalf("could not read config! %s", err)
	}

	// fmt.Println(config)
	fmt.Println(md.Undecoded())

	hub := Hub{
		channels: make(map[string]channels.Channel),
		Env:      NewEnvironment(config),
	}

	for name, cfg := range config.Channels {
		channel, _ := channels.Build(cfg.Class, cfg)

		hub.channels[name] = channel
	}

	return &hub
}

// Run kicks the whole thing off. It should never exit.
func (hub *Hub) Run() {
	events := make(chan event.Event)

	for name, channel := range hub.channels {
		log.Printf("starting channel %s\n", name)
		go channel.Run(events)
	}

	for {
		select {
		case event := <-events:
			fmt.Println(event)
		}
	}
}
