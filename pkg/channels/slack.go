package channels

import (
	"github.com/mmcclimon/synergy-go/internal/slack"
	"github.com/mmcclimon/synergy-go/pkg/event"
)

// SlackChannel is a slack channel.
type SlackChannel struct {
	client *slack.Client
}

// NewSlack gives you a new slack channel
func NewSlack() *SlackChannel {
	channel := SlackChannel{
		client: slack.NewClient(),
	}

	return &channel
}

// Run is the run loop.
func (c *SlackChannel) Run(events chan<- event.Event) {
	rawEvents := make(chan slack.Message)

	go c.client.Run(rawEvents)

	// grab our raw events off of the wire, create synergy events, and pipe them
	// back through to the hub to be handled
	for {
		select {
		case slackEvent := <-rawEvents:
			synergyEvent := event.Event{
				Handled:  false,
				IsPublic: true,
				Type:     "message",
				Text:     slackEvent.Text,
			}

			events <- synergyEvent
		}
	}
}
