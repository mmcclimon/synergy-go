package channels

import (
	"fmt"

	"github.com/mmcclimon/synergy-go/pkg/user"
)

// Event represents a thing to be handled
type Event struct {
	Type                string
	FromUser            *user.User
	Text                string
	IsPublic            bool
	WasTargeted         bool
	FromAddress         string
	ConversationAddress string
	FromChannelName     string
	Handled             bool
}

// Reply sends text to the channel from whence it came.
func (e *Event) Reply(text string) {
	fmt.Printf("would send %s to %s on %s\n", text, e.FromUser.Username, e.FromChannelName)
}
