package hub

import (
	"fmt"

	"github.com/mmcclimon/synergy-go/pkg/channels"
)

type Hub struct {
	name     string
	channels map[string]channels.Channel
}

func NewHub(name string) *Hub {
	hub := Hub{name: name}
	hub.channels = make(map[string]channels.Channel)
	hub.channels["slack"] = channels.NewSlack( /* ... */ )
	return &hub
}

func (hub *Hub) Run() {
	fmt.Printf("running stuff from hub named %s\n", hub.name)
	hub.channels["slack"].Run()
}
