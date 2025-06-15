package main

import (
	"context"
	"log"

	"github.com/jonwilberg/stream-finder/internal/repos/imdb"
	"github.com/jonwilberg/stream-finder/internal/titles"
)

func main() {
	ctx := context.Background()
	imdb.NewIMDBRepository().GetTitles()
	if err := titles.UpdateTitles(ctx); err != nil {
		log.Fatalf("Error updating titles: %v", err)
	}
}
