package titles

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
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
	firestoreClient, err := firestore_repo.NewFirestoreClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create firestore client: %w", err)
	}

	if err := upsertImdbTitles(ctx, firestoreClient); err != nil {
		return fmt.Errorf("failed to upsert imdb titles: %w", err)
	}

	if err := updateNetflixTitles(ctx, firestoreClient); err != nil {
		return fmt.Errorf("failed to update netflix titles: %w", err)
	}

	return nil
}

func upsertImdbTitles(ctx context.Context, firestoreClient *firestore.Client) error {
	imdbRepo := imdb.NewIMDBRepository()

	imdbTitles, err := imdbRepo.GetTitles()

	if err != nil {
		return fmt.Errorf("failed to fetch imdb titles: %w", err)
	}

	documents := make([]firestore_repo.Document, 0, len(imdbTitles))
	for _, title := range imdbTitles {
		documents = append(documents, firestore_repo.Document{
			ID: title.ID,
			Data: Title{
				Title:         title.Title,
				Year:          title.Year,
				UpdatedAt:     time.Now(),
				OriginalTitle: title.OriginalTitle,
				IsAdult:       title.IsAdult,
				Genres:        title.Genres,
				TitleType:     title.TitleType,
			},
		})
	}

	slog.Info("Writing new titles to firestore", "count", len(documents))
	return firestore_repo.BulkWrite(ctx, firestoreClient, "imdb_titles", documents)
}

func updateNetflixTitles(ctx context.Context, firestoreClient *firestore.Client) error {

	newNetflixTitles, err := fetchNewNetflixTitles()
	if err != nil {
		return fmt.Errorf("failed to fetch new netflix titles: %w", err)
	}

	if err := deleteRemovedTitles(ctx, firestoreClient, newNetflixTitles); err != nil {
		return fmt.Errorf("failed to delete removed titles: %w", err)
	}

	return writeNewTitles(ctx, firestoreClient, newNetflixTitles)
}

func fetchNewNetflixTitles() ([]netflix.NetflixTitle, error) {
	netflixRepo := netflix.NewNetflixRepository()
	titles, err := netflixRepo.GetTitles()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch titles from Netflix: %w", err)
	}
	return titles, nil
}

func deleteRemovedTitles(ctx context.Context, client *firestore.Client, newTitles []netflix.NetflixTitle) error {
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

func writeNewTitles(ctx context.Context, client *firestore.Client, titles []netflix.NetflixTitle) error {
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
