package synergy

import (
	"log"

	"github.com/BurntSushi/toml"
)

// FromFile gives you a new hub based on a .toml file. This is the main entry
// point for our package: call this, then call Run() on the thing it gives
// you.
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
