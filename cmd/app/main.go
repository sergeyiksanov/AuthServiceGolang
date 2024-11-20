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

	go func() {
		if err := a.RunMetrics(); err != nil {
			log.Fatalf("Failed to run metrics: %v", err)
		}
		log.Print(err)
	}()

	err = a.Run()
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
