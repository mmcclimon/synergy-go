# synergy.go

This is a port of [Synergy](https://github.com/rjbs/Synergy) to Go.
It's not useful yet, and probably won't be.

## design

You can kick the whole thing off by just running it, which effectly says 
`synergy.FromFile("config.toml").Run()`. Here's the basic design, which is
just a clone of Synergy Prime, altered where needed to be sort of idiomatic in
go.


```
Hub
+---Environment
|   +---UserDirectory
|
+---Channels
+---Reactors
```

The hub is the coordinating object. It has an environment, which is
effectively, global config (including a registry of users).

_Channels_ (not go channels; unfortunate naming clash) represent a means by
which we can send and receive messages. Right now there's only a Slack
channel, but Synergy Prime has a Console channel (for using interactively at
a terminal) and a Twilio channel (for text messaging). The hub kicks off all
the channels, which then pipe Events back through to the hub.

_Reactors_ represent things which respond to events. When the hub receives an
event from a channel, it queries all of its reactors to see if they want to
respond. Reactors might have state, or they might not, but they definitely
have a set of predicates (to see if they should respond to a given event) and
a set of handlers (to actually handle them). A reactor handler is passed the
event, and then expected to call `event.Reply()` so that the user actually
sees a thing.
