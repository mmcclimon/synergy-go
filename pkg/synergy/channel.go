package synergy

import (
	"errors"
	"log"
)

// Channel is a thing on which we can send and receive messages
type Channel interface {
	Name() string
	Run(chan<- Event)
	SendMessage(string, string)
}

// ChannelBuilder is just a convenience type for a .New function
type ChannelBuilder = func(string, ComponentConfig, *Environment) Channel

var channelRegistry = make(map[string]ChannelBuilder)

// RegisterChannel lets you register a builder for this channel
func RegisterChannel(wellKnown string, f ChannelBuilder) {
	log.Printf("registering reactor %s\n", wellKnown)
	channelRegistry[wellKnown] = f
}

// BuildChannel gives you a channel based on a well-known name
func BuildChannel(name string, wellKnown string, cfg ComponentConfig, env *Environment) (Channel, error) {
	builder, ok := channelRegistry[wellKnown]

	if !ok {
		log.Fatalf("unknown reactor name %s", wellKnown)
		return nil, errors.New("unreachable")
	}

	return builder(name, cfg, env), nil
}
