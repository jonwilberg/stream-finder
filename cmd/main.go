package main

import (
	"fmt"
	"log"

	"github.com/jonwilberg/stream-finder/internal/netflix"
)

func main() {
	client := netflix.NewClient()
	repo := netflix.NewRepository(client)

	genreID := "1365" // Action & Adventure
	titles, err := repo.GetGenreTitles(genreID)
	if err != nil {
		log.Fatalf("Error getting titles: %v", err)
	}

	fmt.Printf("Found %d titles in genre %s:\n", len(titles), genreID)
	for _, title := range titles {
		fmt.Println(title)
	}
}
