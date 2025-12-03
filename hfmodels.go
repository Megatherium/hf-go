// Package hfmodels provides a client for searching HuggingFace models
package hfmodels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Megatherium/hf-go/internal/api"
	"github.com/Megatherium/hf-go/internal/models"
)

// Model represents a HuggingFace model
type Model = models.Model

// ListModelsOptions contains options for listing models
type ListModelsOptions = models.ListModelsOptions

// ModelDetails contains detailed model information including files
type ModelDetails struct {
	ID           string    `json:"id"`
	Author       string    `json:"author"`
	Downloads    int       `json:"downloads"`
	Likes        int       `json:"likes"`
	LastModified time.Time `json:"lastModified"`
	PipelineTag  string    `json:"pipeline_tag"`
	LibraryName  string    `json:"library_name"`
	Tags         []string  `json:"tags"`
	Siblings     []Sibling `json:"siblings"`
	CardData     CardData  `json:"cardData"`
	GGUFInfo     *GGUFInfo `json:"gguf"`
}

// Sibling represents a file in the model repository
type Sibling struct {
	RFilename string `json:"rfilename"`
}

// CardData contains model card metadata
type CardData struct {
	ModelName   string      `json:"model_name"`
	ModelType   string      `json:"model_type"`
	BaseModel   interface{} `json:"base_model"` // Can be string or []string
	License     interface{} `json:"license"`    // Can be string or []string
	QuantizedBy string      `json:"quantized_by"`
}

// GetBaseModel returns base_model as a string (first one if array)
func (c CardData) GetBaseModel() string {
	switch v := c.BaseModel.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// GetLicense returns license as a string (first one if array)
func (c CardData) GetLicense() string {
	switch v := c.License.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// GGUFInfo contains GGUF-specific model information
type GGUFInfo struct {
	Total         int64  `json:"total"`
	Architecture  string `json:"architecture"`
	ContextLength int    `json:"context_length"`
}

// Client is a HuggingFace API client
type Client struct {
	client     *api.Client
	httpClient *http.Client
	token      string
}

// NewClient creates a new HuggingFace client
func NewClient(token string) *Client {
	return &Client{
		client:     api.NewClient(token),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}
}

// ListModels fetches models from HuggingFace Hub
func (c *Client) ListModels(opts ListModelsOptions) ([]Model, error) {
	return c.client.ListModels(opts)
}

// GetModelDetails fetches detailed information about a specific model
func (c *Client) GetModelDetails(modelID string) (*ModelDetails, error) {
	url := fmt.Sprintf("https://huggingface.co/api/models/%s", modelID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var details ModelDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

// GetAvailableQuants returns the available quantizations for a GGUF model
func (c *Client) GetAvailableQuants(modelID string) ([]string, error) {
	details, err := c.GetModelDetails(modelID)
	if err != nil {
		return nil, err
	}

	return ExtractQuantsFromSiblings(details.Siblings), nil
}

// ExtractQuantsFromSiblings parses GGUF filenames to extract quantization types
func ExtractQuantsFromSiblings(siblings []Sibling) []string {
	// Match common quant patterns - handles both filename and path formats
	// Patterns: Q4_K_M, IQ4_NL, F16, BF16, TQ1_0, UD-TQ1_0, etc. (case insensitive)
	quantPatterns := []*regexp.Regexp{
		// Standard quants in filename: model-Q4_K_M.gguf or model.Q4_K_M.gguf
		regexp.MustCompile(`(?i)[._-](Q[0-9]+_[A-Z0-9_]+)\.gguf`),
		regexp.MustCompile(`(?i)[._-](IQ[0-9]+_[A-Z0-9_]+)\.gguf`),
		regexp.MustCompile(`(?i)[._-](F16|F32|BF16)\.gguf`),
		// Unsloth style: model-UD-TQ1_0.gguf
		regexp.MustCompile(`(?i)[._-]((?:UD-)?TQ[0-9]+_[0-9]+)\.gguf`),
		// Split files pattern: model-Q4_K_M-00001-of-00005.gguf
		regexp.MustCompile(`(?i)[._-](Q[0-9]+_[A-Z0-9_]+)-\d+-of-\d+\.gguf`),
		regexp.MustCompile(`(?i)[._-](IQ[0-9]+_[A-Z0-9_]+)-\d+-of-\d+\.gguf`),
		regexp.MustCompile(`(?i)[._-](F16|F32|BF16)-\d+-of-\d+\.gguf`),
	}

	// Pattern for directory-based quants: BF16/model-BF16-00001-of-00005.gguf
	dirQuantPattern := regexp.MustCompile(`(?i)^([A-Z0-9_]+)/`)

	seen := make(map[string]bool)
	var quants []string

	for _, s := range siblings {
		filename := s.RFilename
		if !strings.HasSuffix(strings.ToLower(filename), ".gguf") {
			continue
		}

		var quant string

		// First check if it's in a quant-named directory
		if dirMatch := dirQuantPattern.FindStringSubmatch(filename); len(dirMatch) > 1 {
			potentialQuant := dirMatch[1]
			// Verify it looks like a quant name (not just any directory)
			if isQuantName(potentialQuant) {
				quant = potentialQuant
			}
		}

		// If not found in directory, try filename patterns
		if quant == "" {
			for _, pattern := range quantPatterns {
				matches := pattern.FindStringSubmatch(filename)
				if len(matches) > 1 {
					quant = strings.ToUpper(matches[1])
					break
				}
			}
		}

		if quant != "" && !seen[quant] {
			seen[quant] = true
			quants = append(quants, quant)
		}
	}

	return quants
}

// isQuantName checks if a string looks like a quantization name
func isQuantName(s string) bool {
	quantPrefixes := []string{"Q", "IQ", "F16", "F32", "BF16", "TQ"}
	upper := strings.ToUpper(s)
	for _, prefix := range quantPrefixes {
		if strings.HasPrefix(upper, prefix) {
			return true
		}
	}
	return false
}
