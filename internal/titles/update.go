package titles

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jonwilberg/stream-finder/internal/repos/elasticsearch"
	firestore_repo "github.com/jonwilberg/stream-finder/internal/repos/firestore"
	"github.com/jonwilberg/stream-finder/internal/repos/imdb"
	"github.com/jonwilberg/stream-finder/internal/repos/netflix"
)

type Title struct {
	Title         string    `firestore:"title"`
	Year          int       `firestore:"year"`
	UpdatedAt     time.Time `firestore:"updated_at"`
	OriginalTitle string    `firestore:"original_title"`
	IsAdult       bool      `firestore:"is_adult"`
	Genres        []string  `firestore:"genres"`
	TitleType     string    `firestore:"title_type"`
}

func UpdateTitles(ctx context.Context) error {
	elasticsearchClient, err := elasticsearch.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	elasticsearchRepo := elasticsearch.NewRepository(elasticsearchClient)
	if err := elasticsearchRepo.UpdateIndices(ctx); err != nil {
		return fmt.Errorf("failed to update elasticsearch indices: %w", err)
	}

	if err := upsertImdbTitles(ctx, elasticsearchRepo); err != nil {
		return fmt.Errorf("failed to upsert imdb titles: %w", err)
	}

	return nil
}

func upsertImdbTitles(ctx context.Context, elasticsearchRepo *elasticsearch.Repository) error {
	imdbRepo := imdb.NewIMDBRepository()

	imdbTitles, err := imdbRepo.GetTitles()

	if err != nil {
		return fmt.Errorf("failed to fetch imdb titles: %w", err)
	}

	documents := make([]elasticsearch.TitleDocument, 0, len(imdbTitles))
	for _, title := range imdbTitles {
		documents = append(documents, elasticsearch.TitleDocument{
			ID: title.ID,
			Body: elasticsearch.TitleDocumentBody{
				Title:         title.Title,
				Year:          title.Year,
				OriginalTitle: title.OriginalTitle,
				IsAdult:       title.IsAdult,
				Genres:        title.Genres,
				TitleType:     title.TitleType,
			},
		})
	}

	slog.Info("Writing new titles to elasticsearch", "count", len(documents))
	return elasticsearchRepo.BulkIndexTitles(ctx, documents)
}

func FetchNewNetflixTitles() ([]netflix.NetflixTitle, error) {
	netflixRepo := netflix.NewNetflixRepository()
	titles, err := netflixRepo.GetTitles()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch titles from Netflix: %w", err)
	}
	return titles, nil
}

func DeleteRemovedTitles(ctx context.Context, client *firestore.Client, newTitles []netflix.NetflixTitle) error {
	oldTitles, err := firestore_repo.ReadAll(ctx, client, "netflix_titles")
	if err != nil {
		return fmt.Errorf("failed to read existing titles: %w", err)
	}

	newTitleIDs := make(map[string]struct{}, len(newTitles))
	for _, title := range newTitles {
		newTitleIDs[title.ID] = struct{}{}
	}

	removeTitles := make(map[string]struct{}, len(oldTitles))
	for _, oldTitle := range oldTitles {
		if _, exists := newTitleIDs[oldTitle.ID]; !exists {
			removeTitles[oldTitle.ID] = struct{}{}
		}
	}

	if len(removeTitles) > 0 {
		removeIDs := make([]string, 0, len(removeTitles))
		for id := range removeTitles {
			removeIDs = append(removeIDs, id)
		}

		slog.Info("Deleting removed titles from firestore",
			"new_titles", len(newTitles),
			"old_titles", len(oldTitles),
			"removed_titles", len(removeTitles),
		)

		if err := firestore_repo.BulkDelete(ctx, client, "netflix_titles", removeIDs); err != nil {
			return fmt.Errorf("failed to delete removed titles: %w", err)
		}
	} else {
		slog.Info("No removed titles found")
	}

	return nil
}

func WriteNewTitles(ctx context.Context, client *firestore.Client, titles []netflix.NetflixTitle) error {
	documents := make([]firestore_repo.Document, 0, len(titles))
	for _, title := range titles {
		documents = append(documents, firestore_repo.Document{
			ID: title.ID,
			Data: Title{
				Title:     title.Title,
				Year:      title.Year,
				UpdatedAt: time.Now(),
			},
		})
	}

	slog.Info("Writing new titles to firestore", "count", len(documents))
	return firestore_repo.BulkWrite(ctx, client, "netflix_titles", documents)
}
