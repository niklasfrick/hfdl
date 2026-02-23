package models

// FileMetadata represents a file entry from the Hugging Face API
type FileMetadata struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}
