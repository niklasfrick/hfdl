package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"hfdl/internal/api"
	"hfdl/internal/downloader"
)

type DownloadConfig struct {
	ModelID     string
	OutputDir   string
	FileToFetch string
	Token       string
}

func NewDownloadConfig(modelID, outputDir, fileToFetch, token string) *DownloadConfig {
	return &DownloadConfig{
		ModelID:     modelID,
		OutputDir:   outputDir,
		FileToFetch: fileToFetch,
		Token:       token,
	}
}

func (c *DownloadConfig) Run() error {
	safeModelID := sanitizeFilename(c.ModelID)
	modelDir := filepath.Join(c.OutputDir, safeModelID)

	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create model directory: %w", err)
	}

	fmt.Printf("📁 Model directory: %s\n", modelDir)

	apiClient := api.NewClient(c.Token)
	dl := downloader.NewDownloader(apiClient, modelDir)

	if c.FileToFetch != "" {
		if err := dl.DownloadFile(c.ModelID, c.FileToFetch); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}
		fmt.Printf("✅ Successfully downloaded to %s\n", filepath.Join(modelDir, c.FileToFetch))
		return nil
	}

	files, err := apiClient.ListFiles(c.ModelID)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files found in repository.")
		return nil
	}

	fmt.Printf("Found %d files. Starting download...\n", len(files))

	downloaded, err := dl.DownloadAll(c.ModelID, files)
	if err != nil {
		return fmt.Errorf("download completed with errors: %w", err)
	}

	fmt.Printf("✅ Download complete: %d/%d files downloaded.\n", downloaded, len(files))
	return nil
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	safe := re.ReplaceAllString(name, "_")
	safe = strings.TrimFunc(safe, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})
	if safe == "" {
		safe = "model"
	}
	return safe
}
