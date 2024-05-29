package dto

type SearchResult struct {
	Content []*PreviewContent         `json:"content"`
	Persons []*PersonPreviewWithPhoto `json:"persons"`
}
