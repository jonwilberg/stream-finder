package netflix

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type NetflixRepository interface {
	GetGenreTitles(genreID string) ([]NetflixTitle, error)
}

type NetflixTitle struct {
	ID    string
	Title string
	Year  int
}

type netflixRepository struct {
	client *NetflixClient
}

func NewNetflixRepository(client *NetflixClient) NetflixRepository {
	return &netflixRepository{
		client: client,
	}
}

func (r *netflixRepository) GetGenreTitles(genreID string) ([]NetflixTitle, error) {
	body, err := r.client.MakeGenreRequest(genreID)
	if err != nil {
		return nil, fmt.Errorf("failed to make genre request: %w", err)
	}

	videoIDs := extractVideoIDs(body)

	miniModalData, err := r.client.MakeMiniModalRequest(videoIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to make mini modal request: %w", err)
	}

	titles, err := extractTitles(miniModalData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract titles: %w", err)
	}

	return titles, nil
}

func extractVideoIDs(response []byte) []string {
	re := regexp.MustCompile(`Video:(\d+)`)
	matches := re.FindAllStringSubmatch(string(response), -1)

	videoIDs := make([]string, 0, len(matches))
	for _, match := range matches {
		videoIDs = append(videoIDs, match[1])
	}

	return videoIDs
}

func extractTitles(response []byte) ([]NetflixTitle, error) {
	var result struct {
		Data struct {
			UnifiedEntities []struct {
				Title      string `json:"title"`
				VideoID    string `json:"unifiedEntityId"`
				LatestYear int    `json:"latestYear"`
			} `json:"unifiedEntities"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	titles := make([]NetflixTitle, 0, len(result.Data.UnifiedEntities))
	for _, entity := range result.Data.UnifiedEntities {
		titles = append(titles, NetflixTitle{
			ID:    entity.VideoID,
			Title: entity.Title,
			Year:  entity.LatestYear,
		})
	}

	return titles, nil
}
