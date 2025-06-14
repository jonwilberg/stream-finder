package netflix

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractVideoIDs(t *testing.T) {
	sampleData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "netflix_genre_sample.json"))
	if err != nil {
		t.Fatalf("Failed to read sample data: %v", err)
	}

	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "empty response",
			input:    []byte(""),
			expected: []string{},
		},
		{
			name:     "no video IDs",
			input:    []byte(`{"some": "json"}`),
			expected: []string{},
		},
		{
			name:     "single video ID",
			input:    []byte(`{"videos": ["Video:12345"]}`),
			expected: []string{"12345"},
		},
		{
			name:     "multiple video IDs",
			input:    []byte(`{"videos": ["Video:12345", "Video:67890", "Video:11111"]}`),
			expected: []string{"12345", "67890", "11111"},
		},
		{
			name:     "video IDs in different formats",
			input:    []byte(`{"videos": ["Video:12345", "some text", "Video:67890", "more text", "Video:11111"]}`),
			expected: []string{"12345", "67890", "11111"},
		},
		{
			name:     "real API response",
			input:    sampleData,
			expected: []string{"80121192", "81743369"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVideoIDs(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("extractVideoIDs() got %d IDs, want %d", len(got), len(tt.expected))
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("extractVideoIDs()[%d] = %s, want %s", i, id, tt.expected[i])
				}
			}
		})
	}
}
