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

// BuildChannel gives you a channel based on a well-known name
func BuildChannel(name, wellKnown string, cfg ComponentConfig, env *Environment) (Channel, error) {
	switch wellKnown {
	case "SlackChannel":
		return NewSlack(name, cfg, env), nil
	default:
		log.Fatalf("unknown channel name %s", wellKnown)
		return nil, errors.New("unreachable")
	}
}
