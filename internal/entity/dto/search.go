package dto

type SearchResult struct {
	Content []PreviewContentCard     `json:"content"`
	Persons []PersonPreviewWithPhoto `json:"persons"`
}
