package channels

import (
	"errors"
	"log"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/env"
)

// Channel is a thing on which we can send and receive messages
type Channel interface {
	Name() string
	Run(chan<- Event)
	SendMessage(string, string)
}

// Build gives you a channel based on a well-known name
func Build(name, wellKnown string, cfg config.ComponentConfig, env *env.Environment) (Channel, error) {
	switch wellKnown {
	case "SlackChannel":
		return NewSlack(name, cfg, env), nil
	default:
		log.Fatalf("unknown channel name %s", wellKnown)
		return nil, errors.New("unreachable")
	}
}
