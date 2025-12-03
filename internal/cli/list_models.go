package cli

import (
	"fmt"
	"os"

	"github.com/Megatherium/hf-go/internal/api"
	"github.com/Megatherium/hf-go/internal/models"
	"github.com/Megatherium/hf-go/internal/pkg/utils"
	"github.com/spf13/cobra"
)

// ListModelsOptions holds the CLI flags for the list-models command
type ListModelsOptions struct {
	Search       string
	Filter       string
	Author       string
	PipelineTag  string
	LibraryName  string
	Language     string
	Tag          string
	Limit        int
	Sort         string
	Direction    int
	OutputFormat string
	Token        string
}

// NewListModelsCmd creates the list-models command
func NewListModelsCmd() *cobra.Command {
	opts := &ListModelsOptions{}

	cmd := &cobra.Command{
		Use:   "list-models",
		Short: "List models from the Hugging Face Hub",
		Long: `List models from the Hugging Face Hub with various filters and output formats.

Examples:
  # List all models (limited to 20 by default)
  hf-go list-models

  # Search for models with "bert" in their name
  hf-go list-models --search bert

  # Filter by author
  hf-go list-models --author google

  # Filter by pipeline tag
  hf-go list-models --pipeline-tag text-classification

  # JSON output for machine processing
  hf-go list-models --search bert --output-format json

  # Limit results and sort by downloads
  hf-go list-models --limit 10 --sort downloads
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListModels(opts)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&opts.Search, "search", "", "Search for models with this string in their id")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter models by library, task, or tags")
	cmd.Flags().StringVar(&opts.Author, "author", "", "Filter models by author (username or organization)")
	cmd.Flags().StringVar(&opts.PipelineTag, "pipeline-tag", "", "Filter models by pipeline tag (e.g., 'text-generation')")
	cmd.Flags().StringVar(&opts.LibraryName, "library-name", "", "Filter models by library (e.g., 'pytorch', 'tensorflow')")
	cmd.Flags().StringVar(&opts.Language, "language", "", "Filter models by language (e.g., 'en', 'fr')")
	cmd.Flags().StringVar(&opts.Tag, "tag", "", "Filter models by specific tag")
	cmd.Flags().IntVar(&opts.Limit, "limit", 20, "Maximum number of models to return")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort results by field (e.g., 'downloads', 'likes', 'trending_score')")
	cmd.Flags().IntVar(&opts.Direction, "direction", 0, "Sort direction: -1 for descending, 1 for ascending")
	cmd.Flags().StringVar(&opts.OutputFormat, "output-format", "table", "Output format: 'table' or 'json'")
	cmd.Flags().StringVar(&opts.Token, "token", "", "Hugging Face API token (optional, can also use HF_TOKEN env var)")

	return cmd
}

// runListModels executes the list-models command
func runListModels(opts *ListModelsOptions) error {
	// Get token from environment if not provided
	token := opts.Token
	if token == "" {
		token = os.Getenv("HF_TOKEN")
	}

	// Create API client
	client := api.NewClient(token)

	// Build API options
	apiOpts := models.ListModelsOptions{
		Search:      opts.Search,
		Filter:      opts.Filter,
		Author:      opts.Author,
		PipelineTag: opts.PipelineTag,
		LibraryName: opts.LibraryName,
		Language:    opts.Language,
		Tag:         opts.Tag,
		Limit:       opts.Limit,
		Sort:        opts.Sort,
		Direction:   opts.Direction,
		Token:       token,
	}

	// Fetch models
	modelsList, err := client.ListModels(apiOpts)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	// Format output
	switch opts.OutputFormat {
	case "json":
		output, err := utils.FormatJSON(modelsList)
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		fmt.Println(output)
	case "table":
		output := utils.FormatTable(modelsList)
		fmt.Println(output)
	default:
		return fmt.Errorf("unsupported output format: %s (use 'table' or 'json')", opts.OutputFormat)
	}

	return nil
}

// ListModels is a public function that can be used as a library
func ListModels(opts models.ListModelsOptions, format string) (string, error) {
	client := api.NewClient(opts.Token)

	modelsList, err := client.ListModels(opts)
	if err != nil {
		return "", fmt.Errorf("failed to list models: %w", err)
	}

	switch format {
	case "json":
		return utils.FormatJSON(modelsList)
	case "table":
		return utils.FormatTable(modelsList), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}
