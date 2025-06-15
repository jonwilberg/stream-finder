package main

import (
	"fmt"
	"log"

	"github.com/jonwilberg/stream-finder/internal/netflix"
)

func main() {
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

	entityIds := []string{"Video:81588273", "Video:81696513", "Video:81712178"}
	response, err := client.MakeMiniModalRequest(entityIds)
	if err != nil {
		log.Fatalf("Error getting mini modal data: %v", err)
	}

	fmt.Printf("\nMini Modal Response:\n%s\n", string(response))
}
