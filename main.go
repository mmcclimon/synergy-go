package main

import "github.com/mmcclimon/synergy-go/pkg/hub"

func main() {
	hub := hub.FromFile("config.toml")
	hub.Run()
}
