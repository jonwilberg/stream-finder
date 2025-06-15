package titles

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	firestore_repo "github.com/jonwilberg/stream-finder/internal/repos/firestore"
	"github.com/jonwilberg/stream-finder/internal/repos/netflix"
)

type Title struct {
	Title     string    `firestore:"title"`
	Year      int       `firestore:"year"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

func UpdateTitles(ctx context.Context) error {
	newTitles, err := fetchNewTitles()
	if err != nil {
		return fmt.Errorf("failed to fetch new titles: %w", err)
	}

	firestoreClient, err := firestore_repo.NewFirestoreClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create firestore client: %w", err)
	}

	if err := deleteRemovedTitles(ctx, firestoreClient, newTitles); err != nil {
		return fmt.Errorf("failed to delete removed titles: %w", err)
	}

	return writeNewTitles(ctx, firestoreClient, newTitles)
}

func fetchNewTitles() ([]netflix.NetflixTitle, error) {
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

	return firestore_repo.BulkWrite(ctx, client, "netflix_titles", documents)
}
