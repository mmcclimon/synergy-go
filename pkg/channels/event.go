package channels

import (
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
	FromChannel         Channel
	Handled             bool
}

// Reply sends text to the channel from whence it came.
func (e *Event) Reply(text string) {
	prefix := ""

	if e.FromUser != nil && e.IsPublic {
		prefix = e.FromUser.Username + ": "
	}

	text = prefix + text

	e.FromChannel.SendMessage(e.ConversationAddress, text)
}
