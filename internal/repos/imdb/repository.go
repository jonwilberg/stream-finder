package imdb

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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
	url := "https://datasets.imdbws.com/title.basics.tsv.gz"
	filepath := filepath.Join(os.TempDir(), "title.basics.tsv")

	resp, err := r.downloadFile(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := r.unzipFile(resp, filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(filepath)

	return r.extractTitles(file)
}

func (r *imdbRepository) downloadFile(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (r *imdbRepository) unzipFile(resp *http.Response, filepath string) (*os.File, error) {
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := io.Copy(file, gzipReader); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to seek to start of file: %w", err)
	}

	return file, nil
}

func (r *imdbRepository) extractTitles(file *os.File) ([]IMDBTitle, error) {
	bufReader := bufio.NewReader(file)
	csvr := csv.NewReader(bufReader)
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
