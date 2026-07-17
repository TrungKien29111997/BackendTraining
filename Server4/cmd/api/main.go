package main

import (
	"sever4/internal/app"
	"sever4/internal/config"
)

func main() {
	// Init Config
	cfg := config.NewConfig()

	// Init Application
	application := app.NewApplication(cfg)

	// Init Server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
