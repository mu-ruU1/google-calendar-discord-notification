package main

import (
	"os"

	"github.com/mu-ruU1/google-calendar-discord-notification/cmd/googlecal"
)

func loadEnv(key string) (value string) {
	value, ok := os.LookupEnv(key)

	if !ok {
		os.Exit(1)
	}

	return
}

func main() {
	googlecal.Gmain()
}
