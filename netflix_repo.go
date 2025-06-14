package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type NetflixRepository interface {
	GetGenreTitles(genreID int, offset int) ([]string, error)
}

type netflixRepository struct {
	netflixID       string
	netflixSecureID string
	client          *http.Client
}

func NewNetflixRepository() NetflixRepository {
	return &netflixRepository{
		netflixID:       os.Getenv("NETFLIX_ID"),
		netflixSecureID: os.Getenv("NETFLIX_SECURE_ID"),
		client:          &http.Client{},
	}
}

func (r *netflixRepository) makeGenreRequest(genreID int, offset int) ([]byte, error) {
	url := "https://www.netflix.com/nq/website/memberapi/release/pathEvaluator?original_path=%2Fshakti%2Fmre%2FpathEvaluator"

	formBody := &bytes.Buffer{}
	writer := multipart.NewWriter(formBody)

	part, err := writer.CreateFormField("path")
	if err != nil {
		return nil, fmt.Errorf("error creating form field: %v", err)
	}

	pathStr := fmt.Sprintf(`["genres",%d,"su",{"from":%d,"to":%d},"reference",["availability","episodeCount","queue","summary"]]`,
		genreID, offset, offset+10)
	if _, err := part.Write([]byte(pathStr)); err != nil {
		return nil, fmt.Errorf("error writing form field: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, formBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Host", "netflix.com")
	req.Header.Set("Content-Length", strconv.Itoa(formBody.Len()))
	cookie := fmt.Sprintf("SecureNetflixId=%s; NetflixId=%s",
		r.netflixSecureID,
		r.netflixID)
	req.Header.Set("Cookie", cookie)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func extractVideoIDs(response []byte) []string {
	re := regexp.MustCompile(`Video:(\d+)`)
	matches := re.FindAllStringSubmatch(string(response), -1)

	videoIDs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			videoIDs = append(videoIDs, match[1])
		}
	}

	return videoIDs
}

func (r *netflixRepository) GetGenreTitles(genreID int, offset int) ([]string, error) {
	response, err := r.makeGenreRequest(genreID, offset)
	if err != nil {
		return nil, err
	}

	return extractVideoIDs(response), nil
}

func main() {
	repo := NewNetflixRepository()

	videoIDs, err := repo.GetGenreTitles(34399, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d video IDs:\n", len(videoIDs))
	for _, id := range videoIDs {
		fmt.Printf("Video:%s\n", id)
	}
}
