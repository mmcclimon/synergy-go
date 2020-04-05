package event

// Event represents a thing to be handled
type Event struct {
	Handled  bool
	IsPublic bool
	Type     string
	Text     string
}
