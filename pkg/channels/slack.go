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
	channel := SlackChannel{}
	client := slack.NewClient()
	channel.client = client
	return &channel
}

func (c *SlackChannel) Run() {
	c.client.Connect()
}
