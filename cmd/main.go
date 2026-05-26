package main

import (
	"context"
	"log"
)

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// agents起動
	if err := app.AgentRuntime.Run(ctx); err != nil {
		log.Fatal(err)
	}

	// presentation層追加後にUIを起動
	if err := app.UI.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
