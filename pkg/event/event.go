package event

import "github.com/mmcclimon/synergy-go/pkg/user"

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
