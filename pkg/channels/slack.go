package channels

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mmcclimon/synergy-go/internal/config"
	"github.com/mmcclimon/synergy-go/internal/slack"
	"github.com/mmcclimon/synergy-go/pkg/env"
)

// SlackChannel is a slack channel.
type SlackChannel struct {
	name   string
	client *slack.Client
	env    *env.Environment
}

// NewSlack gives you a new slack channel
func NewSlack(name string, cfg config.ComponentConfig, env *env.Environment) *SlackChannel {
	channel := SlackChannel{
		name:   name,
		env:    env,
		client: slack.NewClient(cfg.APIToken),
	}

	return &channel
}

// Name returns the name of this channel
func (c *SlackChannel) Name() string { return c.name }

// Run is the run loop.
func (c *SlackChannel) Run(events chan<- Event) {
	rawEvents := make(chan slack.Message)

	go c.client.Run(rawEvents)

	// grab our raw events off of the wire, create synergy events, and pipe them
	// back through to the hub to be handled
	for {
		select {
		case slackEvent := <-rawEvents:
			fmt.Printf("%#v\n", slackEvent)

			if slackEvent.Type == "" && slackEvent.ReplyTo != 0 {
				if !slackEvent.OK {
					log.Printf("failed to send response: %#v", slackEvent)
				}
				continue
			}

			if slackEvent.BotID != "" {
				continue
			}

			synergyEvent, ok := c.synergyEventFrom(slackEvent)

			if !ok {
				log.Printf(
					"couldn't convert a %s message to channel %s, dropping it",
					slackEvent.Type,
					slackEvent.Channel,
				)
				continue
			}

			events <- *synergyEvent
		}
	}
}

func (c *SlackChannel) synergyEventFrom(slackEvent slack.Message) (*Event, bool) {
	ok := strings.HasPrefix(slackEvent.Channel, "G")

	if !ok {
		privateAddr, ok := c.client.DMChannelForAddress(slackEvent.User)
		if !ok || privateAddr == slack.CannotDMBot {
			return nil, false
		}
	}

	user := c.env.UserDirectory.UserByChannelAndAddress(c.name, slackEvent.User)

	text := c.decodeSlackFormatting(slackEvent.Text)

	targeted := false

	me := c.client.OurName
	// in perl, this uses a lookahead, which you cannot do in go. alas.
	targetedRegex := regexp.MustCompile(`(?i)` + `^@?` + me + `:?\s*`)

	if targetedRegex.MatchString(text) {
		text = targetedRegex.ReplaceAllString(text, "")
		targeted = true
	}

	// everything in DM is targeted
	if slackEvent.Channel[0] == 'D' {
		targeted = true
	}

	// only public channels are public
	isPublic := slackEvent.Channel[0] == 'C'

	synergyEvent := Event{
		Type:                "message",
		Text:                text,
		WasTargeted:         targeted,
		IsPublic:            isPublic,
		FromChannel:         c,
		FromAddress:         slackEvent.User,
		FromUser:            user,
		ConversationAddress: slackEvent.Channel,
		Handled:             false,
		// TransportData?
	}

	return &synergyEvent, true
}

func (c *SlackChannel) decodeSlackFormatting(text string) string {
	// usernames: <@U12345>
	text = regexp.MustCompile(`<@(U[A-Z0-9]+)>`).ReplaceAllStringFunc(text, func(match string) string {
		match = strings.Trim(match, "<>")
		match = match[1:]
		return "@" + c.client.UsernameFor(match)
	})

	// Channels <#C123ABC|bottest>
	text = regexp.MustCompile(`<#[CD](?:[A-Z0-9]+)\|(.*?)>`).ReplaceAllString(text, "#$1")

	// mailto: <mailto:foo@bar.com|foo@bar.com> (no surrounding brackets)
	text = regexp.MustCompile(`<mailto:\S+?\|([^>]+)>`).ReplaceAllString(text, "$1")

	// "helpful" url formatting:  <https://example.com|example.com>; keep what
	// user actually typed
	text = regexp.MustCompile(`<([^>]+)>`).ReplaceAllStringFunc(text, func(match string) string {
		match = strings.Trim(match, "<>")
		return regexp.MustCompile(`^.*\|`).ReplaceAllString(match, "")
	})

	// Anything with < and > around it is probably a URL at this point so remove
	// those
	text = strings.ReplaceAll(text, "<", "")
	text = strings.ReplaceAll(text, ">", "")

	// re-encode html
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")

	return text
}

// SendMessage sends a message
func (c *SlackChannel) SendMessage(target, text string) {
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "&", "&amp;")

	c.client.SendMessage(target, text)
}
