package netflix

import (
	"fmt"
	"io"
	"regexp"
)

// Repository defines the interface for Netflix API operations
type Repository interface {
	GetGenreTitles(genreID string) ([]string, error)
}

// repository implements the Repository interface
type repository struct {
	client *Client
}

// NewRepository creates a new Netflix repository instance
func NewRepository(client *Client) Repository {
	return &repository{
		client: client,
	}
}

// GetGenreTitles retrieves all video IDs for a given genre
func (r *repository) GetGenreTitles(genreID string) ([]string, error) {
	resp, err := r.client.MakeGenreRequest(genreID)
	if err != nil {
		return nil, fmt.Errorf("failed to make genre request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return extractVideoIDs(body), nil
}

// extractVideoIDs parses the API response and extracts video IDs
func extractVideoIDs(response []byte) []string {
	re := regexp.MustCompile(`Video:(\d+)`)
	matches := re.FindAllStringSubmatch(string(response), -1)

	videoIDs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			videoIDs = append(videoIDs, match[1])
		}
	}

	return videoIDs
}
