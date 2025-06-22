package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Document struct {
	ID   string
	Data any
}

func BulkWrite(ctx context.Context, client *firestore.Client, collectionName string, documents []Document) error {
	collection := client.Collection(collectionName)
	bulkWriter := client.BulkWriter(ctx)

	for _, doc := range documents {
		_, err := bulkWriter.Set(collection.Doc(doc.ID), doc.Data)
		if err != nil {
			return fmt.Errorf("failed to create bulk writer job: %w", err)
		}
	}

	bulkWriter.End()
	return nil
}

func BulkDelete(ctx context.Context, client *firestore.Client, collectionName string, documentIDs []string) error {
	collection := client.Collection(collectionName)
	bulkWriter := client.BulkWriter(ctx)

	jobs := make([]*firestore.BulkWriterJob, 0, len(documentIDs))
	for _, id := range documentIDs {
		job, err := bulkWriter.Delete(collection.Doc(id))
		if err != nil {
			return fmt.Errorf("failed to create delete job: %w", err)
		}
		jobs = append(jobs, job)
	}

	bulkWriter.End()
	for _, job := range jobs {
		if _, err := job.Results(); err != nil {
			return fmt.Errorf("firestore delete job failed: %w", err)
		}
	}
	return nil
}

func ReadAll(ctx context.Context, client *firestore.Client, collectionName string) ([]Document, error) {
	collection := client.Collection(collectionName)
	iter := collection.Documents(ctx)
	defer iter.Stop()

	var documents []Document
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}

		documents = append(documents, Document{
			ID:   doc.Ref.ID,
			Data: doc.Data(),
		})
	}

	return documents, nil
}
