package titles

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jonwilberg/stream-finder/internal/repos/elasticsearch"
)

func SearchTitles(ctx context.Context, elasticsearchRepo *elasticsearch.Repository, query string, limit int) ([]elasticsearch.TitleDocument, error) {
	searchQuery := map[string]any{
		"size": limit,
		"query": map[string]any{
			"match_phrase": map[string]any{
				"title": query,
			},
		},
		"sort": []map[string]any{
			{"_score": map[string]any{"order": "desc"}},
		},
	}

	responseBytes, err := elasticsearchRepo.Search(ctx, "titles", searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search titles: %w", err)
	}

	var response elasticsearch.SearchResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	results := make([]elasticsearch.TitleDocument, 0, len(response.Hits.Hits))
	for _, hit := range response.Hits.Hits {
		results = append(results, elasticsearch.TitleDocument{
			ID:   hit.ID,
			Body: hit.Source,
		})
	}

	return results, nil
}
