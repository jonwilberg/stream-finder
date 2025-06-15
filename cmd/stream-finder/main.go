package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jonwilberg/stream-finder/internal/firestore"
	"github.com/jonwilberg/stream-finder/internal/netflix"
)

func main() {
	ctx := context.Background()
	firestoreClient, err := firestore.NewFirestoreClient(ctx, "stream-finder-463006")
	if err != nil {
		log.Fatalf("Error creating firestore client: %v", err)
	}

	firestoreClient.Collection("netflix_titles").Doc("Video:81588273").Set(ctx, map[string]any{
		"title": "The Dark Knight",
		"year":  2008,
	})

	client := netflix.NewClient()
	repo := netflix.NewNetflixRepository(client)

	genreID := "1365" // Action & Adventure
	titles, err := repo.GetGenreTitles(genreID)
	if err != nil {
		log.Fatalf("Error getting titles: %v", err)
	}

	fmt.Printf("Found %d titles in genre %s:\n", len(titles), genreID)
	for _, title := range titles {
		fmt.Printf("%s (%d)\n", title.Title, title.Year)
	}
}
