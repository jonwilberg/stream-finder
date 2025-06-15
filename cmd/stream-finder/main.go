package main

import (
	"context"
	"log"

	"github.com/jonwilberg/stream-finder/internal/titles"
)

func main() {
	ctx := context.Background()
	if err := titles.UpdateTitles(ctx); err != nil {
		log.Fatalf("Error updating titles: %v", err)
	}
}
