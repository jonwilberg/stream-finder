package elasticsearch

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	*elasticsearch.Client
}

func NewClient() (*Client, error) {
	password, err := getElasticsearchPassword()
	if err != nil {
		return nil, err
	}

	cfg := elasticsearch.Config{
		Addresses: []string{getElasticsearchURL()},
		Username:  getElasticsearchUsername(),
		Password:  password,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 30 * time.Second,
		},
		MaxRetries: 3,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &Client{Client: client}, nil
}

func getElasticsearchURL() string {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://localhost:9200"
	}
	return url
}

func getElasticsearchPassword() (string, error) {
	password := os.Getenv("ELASTICSEARCH_PASSWORD")
	if password == "" {
		return "", fmt.Errorf("ELASTICSEARCH_PASSWORD environment variable is required")
	}
	return password, nil
}

func getElasticsearchUsername() string {
	username := os.Getenv("ELASTICSEARCH_USERNAME")
	if username == "" {
		username = "elastic"
	}
	return username
}
