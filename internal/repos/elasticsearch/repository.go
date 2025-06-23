package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/jonwilberg/stream-finder/pkg/logging"
)

type TitleDocument struct {
	ID   string
	Body TitleDocumentBody
}

type TitleDocumentBody struct {
	TitleType     string   `json:"title_type"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title"`
	IsAdult       bool     `json:"is_adult"`
	Year          int      `json:"year"`
	Genres        []string `json:"genres"`
}

type Repository struct {
	client *Client
}

func NewRepository(client *Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) BulkIndexTitles(ctx context.Context, titleDocs []TitleDocument) error {
	bulkIndexerConfig := esutil.BulkIndexerConfig{
		Index:         "titles",
		Client:        r.client,
		NumWorkers:    10,
		FlushBytes:    5_000_000,
		FlushInterval: 30 * time.Second,
	}

	bi, err := esutil.NewBulkIndexer(bulkIndexerConfig)
	if err != nil {
		return fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	bar := logging.NewProgressBar("Indexing titles to Elasticsearch", len(titleDocs))

	for _, doc := range titleDocs {
		docJSON, err := json.Marshal(doc.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: doc.ID,
			Body:       bytes.NewReader(docJSON),
		})
		bar.Add(1)
	}

	bar.Finish()

	if err := bi.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}

	return nil
}

func (r *Repository) EnsureIndexExists(ctx context.Context, indexName string, mappingJSON string) error {
	exists, err := r.client.Indices.Exists([]string{indexName}, r.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}

	if exists.StatusCode == 200 {
		return nil
	}

	resp, err := r.client.Indices.Create(
		indexName,
		r.client.Indices.Create.WithContext(ctx),
		r.client.Indices.Create.WithBody(bytes.NewReader([]byte(mappingJSON))),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create index: %s", string(bodyBytes))
	}

	return nil
}

func (r *Repository) UpdateIndices(ctx context.Context) error {
	schemasDir := "internal/repos/elasticsearch/schemas"

	entries, err := os.ReadDir(schemasDir)
	if err != nil {
		return fmt.Errorf("failed to read schemas directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		indexName := strings.TrimSuffix(entry.Name(), ".json")
		schemaPath := filepath.Join(schemasDir, entry.Name())

		mappingJSON, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
		}

		if err := r.EnsureIndexExists(ctx, indexName, string(mappingJSON)); err != nil {
			return fmt.Errorf("failed to ensure index %s exists: %w", indexName, err)
		}
	}

	return nil
}
