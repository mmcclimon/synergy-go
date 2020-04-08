package channels

import (
	"errors"
	"log"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/pkg/env"
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// Channel is a thing on which we can send and receive messages
type Channel interface {
	Run(chan<- event.Event)
}

// Build gives you a channel based on a well-known name
func Build(name, wellKnown string, cfg config.ChannelConfig, env *env.Environment) (Channel, error) {
	switch wellKnown {
	case "SlackChannel":
		return NewSlack(name, cfg, env), nil
	default:
		log.Fatalf("unknown channel name %s", wellKnown)
		return nil, errors.New("unreachable")
	}
}
