package main

import (
	"context"
	"log"

	"github.com/kakkky/hakoniwa/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// agents起動
	if err := app.Runtime.Run(ctx); err != nil {
		log.Fatal(err)
	}

	// presentation層追加後にUIを起動
	// if err := app.UI.Run(ctx); err != nil {
	// 	log.Fatal(err)
	// }
}
