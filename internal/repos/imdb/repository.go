package imdb

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/jonwilberg/stream-finder/pkg/logging"
	"github.com/jszwec/csvutil"
)

type IMDBRepository interface {
	GetTitles() ([]IMDBTitle, error)
}

type GenreList []string

func (g *GenreList) UnmarshalCSV(data []byte) error {
	s := string(data)
	if s == `\N` { // IMDb uses "\N" for null
		*g = nil
		return nil
	}
	*g = strings.Split(s, ",")
	return nil
}

type IMDBTitle struct {
	ID            string    `csv:"tconst"`
	TitleType     string    `csv:"titleType"`
	Title         string    `csv:"primaryTitle"`
	OriginalTitle string    `csv:"originalTitle"`
	IsAdult       bool      `csv:"isAdult"`
	Year          int       `csv:"startYear"`
	Genres        GenreList `csv:"genres"`
}

type imdbRepository struct{}

func NewIMDBRepository() IMDBRepository {
	return &imdbRepository{}
}

func (r *imdbRepository) GetTitles() ([]IMDBTitle, error) {
	file, err := os.Open("...")

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	csvr := csv.NewReader(reader)
	csvr.Comma = '\t'
	csvr.FieldsPerRecord = 9
	csvr.ReuseRecord = true

	dec, err := csvutil.NewDecoder(csvr)
	if err != nil {
		return nil, fmt.Errorf("failed to create csv decoder: %w", err)
	}

	titles := make([]IMDBTitle, 0, 12_000_000)
	failed := 0
	bar := logging.NewProgressBar("Decoding IMDb titles")
	for {
		var t IMDBTitle
		if err := dec.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			failed++
		} else {
			titles = append(titles, t)
			bar.Add(1)
		}
	}
	bar.Finish()

	slog.Info("Decoded IMDb titles", "count", len(titles), "failed", failed)
	return titles, nil
}
