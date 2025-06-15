package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	return client, nil
}
