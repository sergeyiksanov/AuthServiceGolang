package main

import (
	"AuthService/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
