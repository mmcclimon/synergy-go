package channels

import (
	"fmt"
	"os"

	"github.com/mmcclimon/synergy-go/internal/slack"
)

// SlackChannel is a slack channel.
type SlackChannel struct {
	client *slack.Client
}

// NewSlack gives you a new slack channel
func NewSlack() *SlackChannel {
	channel := SlackChannel{}
	client, err := slack.NewClient()

	if err != nil {
		fmt.Println(err.Error())
		// stupid, but works
		os.Exit(1)
	}

	channel.client = client
	return &channel
}

func (c *SlackChannel) Run() {
	c.client.Connect()
}
