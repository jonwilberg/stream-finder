package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("GC_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	return client, nil
}
