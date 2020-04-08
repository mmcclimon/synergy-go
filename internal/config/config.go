package config

// Config is a config; it has some data in it
type Config struct {
	StateDBFile string `toml:"state_dbfile"`

	Channels map[string]ComponentConfig

	Reactors map[string]ComponentConfig
}

// ComponentConfig is a config for a channel or reactor (still in progress)
type ComponentConfig struct {
	Class string

	// this exists on the slack channel, but wouldn't on (say) a terminal
	// channel. I think that's fine, for now, but it means this must be a union
	// of all possible channel configs, which isn't great.
	APIToken string

	// Maybe...we have "Other" here, and we fill it in from the metadata, as a
	// map[string]interface{} or wevs, and then individual channels can do
	// whatever they like with it.
}
