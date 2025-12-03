package models

import "time"

// Model represents a Hugging Face model with its metadata
type Model struct {
	ID            string    `json:"id"`
	Author        string    `json:"author"`
	Downloads     int       `json:"downloads"`
	Likes         int       `json:"likes"`
	LastModified  time.Time `json:"lastModified"`
	LibraryName   string    `json:"library_name,omitempty"`
	PipelineTag   string    `json:"pipeline_tag,omitempty"`
	Private       bool      `json:"private"`
	Gated         bool      `json:"gated,omitempty"`
	TrendingScore float64   `json:"trending_score,omitempty"`
}

// ListModelsOptions contains parameters for filtering and sorting models
type ListModelsOptions struct {
	Search      string
	Filter      string
	Author      string
	PipelineTag string
	LibraryName string
	Language    string
	Tag         string
	Limit       int
	Sort        string
	Direction   int
	Token       string
}
