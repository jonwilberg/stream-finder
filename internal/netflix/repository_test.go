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

func TestExtractTitles(t *testing.T) {
	sampleData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "netflix_mini_modal_sample.json"))
	if err != nil {
		t.Fatalf("Failed to read sample data: %v", err)
	}

	tests := []struct {
		name     string
		input    []byte
		expected []NetflixTitle
		wantErr  bool
	}{
		{
			name:     "empty response",
			input:    []byte(""),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{"invalid": json}`),
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "real API response",
			input: sampleData,
			expected: []NetflixTitle{
				{
					ID:    "Video:81588273",
					Title: "A Deadly American Marriage",
					Year:  2025,
				},
				{
					ID:    "Video:81696513",
					Title: "The Beekeeper",
					Year:  2023,
				},
				{
					ID:    "Video:81712178",
					Title: "Titan: The OceanGate Submersible Disaster",
					Year:  2025,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTitles(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTitles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("extractTitles() got %d titles, want %d", len(got), len(tt.expected))
				return
			}
			for i, title := range got {
				if title.ID != tt.expected[i].ID {
					t.Errorf("extractTitles()[%d].ID = %v, want %v", i, title.ID, tt.expected[i].ID)
				}
				if title.Title != tt.expected[i].Title {
					t.Errorf("extractTitles()[%d].Title = %v, want %v", i, title.Title, tt.expected[i].Title)
				}
				if title.Year != tt.expected[i].Year {
					t.Errorf("extractTitles()[%d].Year = %v, want %v", i, title.Year, tt.expected[i].Year)
				}
			}
		})
	}
}
