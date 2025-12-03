package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Megatherium/hf-go/internal/models"
)

const (
	DefaultAPIURL = "https://huggingface.co/api/models"
)

// Client represents a Hugging Face API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new Hugging Face API client
func NewClient(token string) *Client {
	return &Client{
		BaseURL: DefaultAPIURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}
}

// apiModel represents the raw model response from the API
type apiModel struct {
	ID            string    `json:"id"`
	Downloads     int       `json:"downloads"`
	Likes         int       `json:"likes"`
	LastModified  time.Time `json:"lastModified"`
	LibraryName   string    `json:"library_name"`
	PipelineTag   string    `json:"pipeline_tag"`
	Private       bool      `json:"private"`
	Gated         interface{} `json:"gated"`
	TrendingScore float64   `json:"trendingScore"`
}

// ListModels fetches models from the Hugging Face Hub based on the provided options
func (c *Client) ListModels(opts models.ListModelsOptions) ([]models.Model, error) {
	// Build query parameters
	params := url.Values{}

	if opts.Search != "" {
		params.Add("search", opts.Search)
	}
	if opts.Filter != "" {
		params.Add("filter", opts.Filter)
	}
	if opts.Author != "" {
		params.Add("author", opts.Author)
	}
	if opts.PipelineTag != "" {
		params.Add("pipeline_tag", opts.PipelineTag)
	}
	if opts.LibraryName != "" {
		params.Add("library", opts.LibraryName)
	}
	if opts.Language != "" {
		params.Add("language", opts.Language)
	}
	if opts.Tag != "" {
		params.Add("tags", opts.Tag)
	}
	if opts.Limit > 0 {
		params.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Sort != "" {
		params.Add("sort", opts.Sort)
	}
	if opts.Direction != 0 {
		params.Add("direction", strconv.Itoa(opts.Direction))
	}

	// Build request URL
	reqURL := c.BaseURL
	if len(params) > 0 {
		reqURL = fmt.Sprintf("%s?%s", c.BaseURL, params.Encode())
	}

	// Create request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header if token is provided
	token := opts.Token
	if token == "" {
		token = c.Token
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiModels []apiModel
	if err := json.Unmarshal(body, &apiModels); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to internal model format
	result := make([]models.Model, len(apiModels))
	for i, am := range apiModels {
		author := ""
		if strings.Contains(am.ID, "/") {
			parts := strings.SplitN(am.ID, "/", 2)
			author = parts[0]
		}

		gated := false
		if am.Gated != nil {
			switch v := am.Gated.(type) {
			case bool:
				gated = v
			case string:
				gated = v != "" && v != "false"
			}
		}

		result[i] = models.Model{
			ID:            am.ID,
			Author:        author,
			Downloads:     am.Downloads,
			Likes:         am.Likes,
			LastModified:  am.LastModified,
			LibraryName:   am.LibraryName,
			PipelineTag:   am.PipelineTag,
			Private:       am.Private,
			Gated:         gated,
			TrendingScore: am.TrendingScore,
		}
	}

	return result, nil
}
