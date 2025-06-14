package netflix

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// Client handles HTTP requests to the Netflix API
type Client struct {
	netflixID       string
	netflixSecureID string
	client          *http.Client
}

// NewClient creates a new Netflix API client
func NewClient() *Client {
	return &Client{
		netflixID:       os.Getenv("NETFLIX_ID"),
		netflixSecureID: os.Getenv("NETFLIX_SECURE_ID"),
		client:          &http.Client{},
	}
}

// MakeGenreRequest sends a request to get titles for a specific genre
func (c *Client) MakeGenreRequest(genreID string) (*http.Response, error) {
	url := "https://www.netflix.com/nq/website/memberapi/release/pathEvaluator?original_path=%2Fshakti%2Fmre%2FpathEvaluator"

	formBody := &bytes.Buffer{}
	writer := multipart.NewWriter(formBody)

	part, err := writer.CreateFormField("path")
	if err != nil {
		return nil, fmt.Errorf("error creating form field: %v", err)
	}

	pathStr := fmt.Sprintf(`["genres",%s,"su",{"from":0,"to":10},"reference",["availability","episodeCount","queue","summary"]]`,
		genreID)
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
		c.netflixSecureID,
		c.netflixID)
	req.Header.Set("Cookie", cookie)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
