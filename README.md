# Hugging Face Models - Go CLI & Library

A Go implementation of the Hugging Face Hub models listing functionality. This project provides both a CLI tool and a reusable library for listing and filtering models from the Hugging Face Hub.

## Features

- List models from Hugging Face Hub with various filters
- Support for multiple output formats (table, JSON)
- CLI tool built with Cobra
- Reusable as a Go module/library
- Filters: search, author, pipeline tag, library, language, tags
- Sorting by downloads, likes, trending score, etc.

## Installation

### As a CLI tool

```bash
cd cmd
go build -o hf-go
./hf-go list-models --help
```

### As a library

```bash
go get github.com/Megatherium/hf-go
```

## Usage

### CLI Examples

```bash
# List all models (limited to 20 by default)
./hf-go list-models

# Search for models with "bert" in their name
./hf-go list-models --search bert

# Filter by author
./hf-go list-models --author google

# Filter by pipeline tag
./hf-go list-models --pipeline-tag text-classification

# JSON output for machine processing
./hf-go list-models --search bert --output-format json

# Limit results and sort by downloads
./hf-go list-models --limit 10 --sort downloads --direction -1

# Combine multiple filters
./hf-go list-models --author openai --pipeline-tag text-generation --limit 5
```

### Library Examples

#### Example 1: Using the API client directly

```go
package main

import (
    "fmt"
    "log"

    "github.com/Megatherium/hf-go/internal/api"
    "github.com/Megatherium/hf-go/internal/models"
)

func main() {
    client := api.NewClient("") // Pass HF token if needed

    modelsList, err := client.ListModels(models.ListModelsOptions{
        Search: "bert",
        Limit:  5,
        Sort:   "downloads",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    for _, model := range modelsList {
        fmt.Printf("%s - Downloads: %d\n", model.ID, model.Downloads)
    }
}
```

#### Example 2: Using the CLI library function

```go
package main

import (
    "fmt"
    "log"

    "github.com/Megatherium/hf-go/internal/cli"
    "github.com/Megatherium/hf-go/internal/models"
)

func main() {
    // Get table formatted output
    output, err := cli.ListModels(models.ListModelsOptions{
        Search: "gpt2",
        Limit:  3,
    }, "table")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(output)

    // Get JSON formatted output
    jsonOutput, err := cli.ListModels(models.ListModelsOptions{
        Author: "google",
        Limit:  2,
    }, "json")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(jsonOutput)
}
```

## API Options

### ListModelsOptions

- `Search` - Search for models with this string in their ID
- `Filter` - Filter by library, task, or tags
- `Author` - Filter by author (username or organization)
- `PipelineTag` - Filter by pipeline tag (e.g., 'text-generation')
- `LibraryName` - Filter by library (e.g., 'pytorch', 'tensorflow')
- `Language` - Filter by language (e.g., 'en', 'fr')
- `Tag` - Filter by specific tag
- `Limit` - Maximum number of models to return
- `Sort` - Sort by field (e.g., 'downloads', 'likes', 'trending_score')
- `Direction` - Sort direction: -1 for descending, 1 for ascending
- `Token` - Hugging Face API token (optional)

## Output Formats

### Table Format (default)

Pretty-printed table with columns:
- Model ID
- Author
- Downloads
- Likes
- Last Modified
- Library
- Task

### JSON Format

Machine-readable JSON array containing model objects with all available metadata.

## Environment Variables

- `HF_TOKEN` - Hugging Face API token (optional, for accessing private models)

## Project Structure

```
go-hf-go/
├── cmd/
│   └── main.go                 # CLI entry point
├── internal/
│   ├── api/
│   │   └── client.go          # API client
│   ├── cli/
│   │   ├── root.go            # Root command
│   │   └── list_models.go     # List models command
│   ├── models/
│   │   └── model.go           # Data models
│   └── pkg/
│       ├── examples/
│       │   └── example.go     # Usage examples
│       └── utils/
│           └── formatters.go  # Output formatters
├── go.mod
└── README.md
```

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework

## License

Apache License 2.0
