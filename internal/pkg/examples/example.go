package main

import (
	"fmt"
	"log"

	"github.com/Megatherium/hf-go/internal/api"
	"github.com/Megatherium/hf-go/internal/cli"
	"github.com/Megatherium/hf-go/internal/models"
	"github.com/Megatherium/hf-go/internal/pkg/utils"
)

func main() {
	// Example 1: Using the API client directly
	fmt.Println("Example 1: Using the API client directly")
	fmt.Println("==========================================")

	client := api.NewClient("")

	modelsList, err := client.ListModels(models.ListModelsOptions{
		Search: "bert",
		Limit:  5,
		Sort:   "downloads",
	})
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	fmt.Printf("Found %d models:\n", len(modelsList))
	for _, model := range modelsList {
		fmt.Printf("- %s (downloads: %d)\n", model.ID, model.Downloads)
	}

	fmt.Println()

	// Example 2: Using the CLI library function with table format
	fmt.Println("Example 2: Using the CLI library function with table format")
	fmt.Println("============================================================")

	output, err := cli.ListModels(models.ListModelsOptions{
		Search: "gpt2",
		Limit:  3,
	}, "table")
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	fmt.Println(output)
	fmt.Println()

	// Example 3: Using the CLI library function with JSON format
	fmt.Println("Example 3: Using the CLI library function with JSON format")
	fmt.Println("===========================================================")

	jsonOutput, err := cli.ListModels(models.ListModelsOptions{
		Author: "google",
		Limit:  2,
	}, "json")
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	fmt.Println(jsonOutput)
	fmt.Println()

	// Example 4: Custom formatting using the formatters
	fmt.Println("Example 4: Custom formatting")
	fmt.Println("=============================")

	models2, err := client.ListModels(models.ListModelsOptions{
		PipelineTag: "text-generation",
		Limit:       3,
	})
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	// Custom processing
	fmt.Println("Text Generation Models:")
	for _, model := range models2 {
		fmt.Printf("  %s - Library: %s\n", model.ID, model.LibraryName)
	}

	fmt.Println()
	fmt.Println("Formatted as JSON:")
	jsonStr, _ := utils.FormatJSON(models2)
	fmt.Println(jsonStr)
}
