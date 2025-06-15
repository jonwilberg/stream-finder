package netflix

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/schollz/progressbar/v3"
)

type NetflixRepository interface {
	GetTitles() ([]NetflixTitle, error)
}

type NetflixTitle struct {
	ID    string
	Title string
	Year  int
}

type netflixRepository struct {
	client *NetflixClient
}

func NewNetflixRepository() NetflixRepository {
	client := NewClient()
	return &netflixRepository{
		client: client,
	}
}

func (r *netflixRepository) GetTitles() ([]NetflixTitle, error) {
	allTitles := []NetflixTitle{}

	genres := []string{
		"1365",  // Movies
		"43040", // TV Shows,
	}

	for _, genreID := range genres {
		titles, err := r.GetGenreTitles(genreID)
		if err != nil {
			return nil, fmt.Errorf("failed to get titles: %w", err)
		}
		allTitles = append(allTitles, titles...)
	}

	return allTitles, nil
}

func (r *netflixRepository) GetGenreTitles(genreID string) ([]NetflixTitle, error) {
	batchSize := 100
	offset := 0
	var allTitles []NetflixTitle

	bar := newProgressBar(fmt.Sprintf("Fetching titles from Netflix genre %s", genreID))

	for {
		titles, err := r.getGenreTitlesBatch(genreID, offset, batchSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get movies titles: %w", err)
		}

		if len(titles) == 0 {
			break
		}

		allTitles = append(allTitles, titles...)
		bar.Add(len(titles))
		offset += batchSize + 1
	}

	bar.Finish()
	slog.Info("Fetched titles from Netflix", "count", len(allTitles))
	return allTitles, nil
}

func (r *netflixRepository) getGenreTitlesBatch(genreID string, offset int, batchSize int) ([]NetflixTitle, error) {
	body, err := r.client.MakeGenreRequest(genreID, offset, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to make genre request: %w", err)
	}

	videoIDs := extractVideoIDs(body)
	if len(videoIDs) == 0 {
		return []NetflixTitle{}, nil
	}

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

	uniqueIDs := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		uniqueIDs["Video:"+match[1]] = struct{}{}
	}

	videoIDs := make([]string, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		videoIDs = append(videoIDs, id)
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

func newProgressBar(description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
}
