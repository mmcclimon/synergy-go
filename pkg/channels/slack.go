package channels

import (
	"github.com/mmcclimon/synergy-go/internal/slack"
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

func (c *SlackChannel) Run() {
	go c.client.Run()
	for {
		// loop forever
	}
}
