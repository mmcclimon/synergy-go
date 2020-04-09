package main

import "github.com/mmcclimon/synergy-go/pkg/synergy"

func main() {
	hub := synergy.FromFile("config.toml")
	hub.Run()
}
