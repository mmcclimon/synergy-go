package main

import "github.com/mmcclimon/synergy-go/pkg/hub"

func main() {
	hub := hub.NewHub("synergy")
	hub.Run()
}
