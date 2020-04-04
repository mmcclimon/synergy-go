package channels

// Channel is a thing on which we can send and receive messages
type Channel interface {
	Run()
}
