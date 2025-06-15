package netflix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type NetflixClient struct {
	netflixID       string
	netflixSecureID string
	client          *http.Client
}

func NewClient() *NetflixClient {
	return &NetflixClient{
		netflixID:       os.Getenv("NETFLIX_ID"),
		netflixSecureID: os.Getenv("NETFLIX_SECURE_ID"),
		client:          &http.Client{},
	}
}

func (c *NetflixClient) MakeGenreRequest(genreID string) ([]byte, error) {
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func (c *NetflixClient) MakeMiniModalRequest(unifiedEntityIds []string) ([]byte, error) {
	url := "https://web.prod.cloud.netflix.com/graphql"

	requestBody := map[string]any{
		"operationName": "MiniModalQuery",
		"variables": map[string]any{
			"videoMerchEnabled":       false,
			"fetchPromoVideoOverride": false,
			"hasPromoVideoOverride":   false,
			"promoVideoId":            0,
			"videoMerchContext":       "BROWSE",
			"isLiveEpisodic":          false,
			"artworkContext":          map[string]any{},
			"textEvidenceUiContext":   "BOB",
			"unifiedEntityIds":        unifiedEntityIds,
		},
		"extensions": map[string]any{
			"persistedQuery": map[string]any{
				"id":      "cea97958-c71c-4c3c-b94c-877fb3c9b89d",
				"version": 102,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "netflix.com")
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonBody)))
	cookie := fmt.Sprintf("SecureNetflixId=%s; NetflixId=%s",
		c.netflixSecureID,
		c.netflixID)
	req.Header.Set("Cookie", cookie)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
