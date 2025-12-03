package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Megatherium/hf-go/internal/models"
)

// FormatTable formats models as a pretty-printed table
func FormatTable(modelsList []models.Model) string {
	if len(modelsList) == 0 {
		return "No models found matching the specified criteria."
	}

	// Headers
	headers := []string{"Model ID", "Author", "Downloads", "Likes", "Last Modified", "Library", "Task"}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	// Prepare rows and update column widths
	rows := make([][]string, len(modelsList))
	for i, model := range modelsList {
		lastModified := "N/A"
		if !model.LastModified.IsZero() {
			lastModified = model.LastModified.Format("2006-01-02")
		}

		library := model.LibraryName
		if library == "" {
			library = "N/A"
		}

		task := model.PipelineTag
		if task == "" {
			task = "N/A"
		}

		row := []string{
			model.ID,
			model.Author,
			formatNumber(model.Downloads),
			formatNumber(model.Likes),
			lastModified,
			library,
			task,
		}
		rows[i] = row

		// Update widths
		for j, cell := range row {
			if len(cell) > widths[j] {
				widths[j] = len(cell)
			}
		}
	}

	// Build table
	var sb strings.Builder

	// Print header separator
	sb.WriteString(buildSeparator(widths))
	sb.WriteString("\n")

	// Print headers
	sb.WriteString("│")
	for i, header := range headers {
		sb.WriteString(fmt.Sprintf(" %-*s │", widths[i], header))
	}
	sb.WriteString("\n")

	// Print header separator
	sb.WriteString(buildSeparator(widths))
	sb.WriteString("\n")

	// Print rows
	for _, row := range rows {
		sb.WriteString("│")
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf(" %-*s │", widths[i], cell))
		}
		sb.WriteString("\n")
	}

	// Print footer separator
	sb.WriteString(buildSeparator(widths))

	return sb.String()
}

// buildSeparator creates a table separator line
func buildSeparator(widths []int) string {
	var sb strings.Builder
	sb.WriteString("├")
	for i, width := range widths {
		sb.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			sb.WriteString("┼")
		}
	}
	sb.WriteString("┤")
	return sb.String()
}

// formatNumber formats an integer with thousands separators
func formatNumber(n int) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}
	return result.String()
}

// FormatJSON formats models as JSON
func FormatJSON(modelsList []models.Model) (string, error) {
	if len(modelsList) == 0 {
		return "[]", nil
	}

	output, err := json.MarshalIndent(modelsList, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(output), nil
}
