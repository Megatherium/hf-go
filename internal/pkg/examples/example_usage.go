// Example usage of the hf-go package as a library
package main

import (
	"fmt"
	"log"

	"github.com/Megatherium/hf-go/internal/cli"
	"github.com/Megatherium/hf-go/internal/models"
)

func main() {
	fmt.Println("Example 1: Search for BERT models with table output")
	fmt.Println("=====================================================")

	output, err := cli.ListModels(models.ListModelsOptions{
		Search: "bert",
		Limit:  5,
	}, "table")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(output)
	fmt.Println()

	fmt.Println("Example 2: Get OpenAI models as JSON")
	fmt.Println("=====================================")

	jsonOutput, err := cli.ListModels(models.ListModelsOptions{
		Author: "openai",
		Limit:  3,
	}, "json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(jsonOutput)
}
