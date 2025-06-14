package main

import (
	"fmt"

	"github.com/jonwilberg/stream-finder/internal/netflix"
)

func main() {
	client := netflix.NewClient()
	repo := netflix.NewRepository(client)

	genreID := "34399"
	titles, err := repo.GetGenreTitles(genreID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d titles:\n", len(titles))
	for _, title := range titles {
		fmt.Printf("Video:%s\n", title)
	}
}
