package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"hfdl/internal/api"
	"hfdl/internal/progress"
	"hfdl/models"
)

type Downloader struct {
	client  *api.Client
	baseDir string
}

func NewDownloader(apiClient *api.Client, outputDir string) *Downloader {
	return &Downloader{
		client:  apiClient,
		baseDir: outputDir,
	}
}

func (d *Downloader) SetToken(token string) {
	if token != "" {
		d.client = api.NewClient(token)
	}
}

func (d *Downloader) DownloadFile(repoID, filePath string) error {
	safePath := filepath.Join(d.baseDir, filePath)
	if !strings.HasPrefix(safePath, filepath.Clean(d.baseDir)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if err := os.MkdirAll(filepath.Dir(safePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if info, err := os.Stat(safePath); err == nil && info.Size() > 0 {
		fmt.Printf("⏭️  Skipping %s (already exists)\n", filePath)
		return nil
	}

	reader, total, err := d.client.DownloadFile(repoID, filePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	outFile, err := os.CreateTemp(filepath.Dir(safePath), ".tmp-"+filepath.Base(safePath))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(outFile.Name())

	bar := progress.NewProgressBar(filePath, total)
	bar.Start()

	written, err := io.Copy(outFile, &progress.ProgressReader{
		Reader: reader,
		Bar:    bar,
	})
	if err != nil {
		outFile.Close()
		return fmt.Errorf("write failed: %w", err)
	}
	outFile.Close()

	if total <= 0 {
		total = written
	}

	if err := os.Rename(outFile.Name(), safePath); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	fmt.Println()
	fmt.Printf("✅ %s (%.2f MB)\n", filePath, float64(total)/1024/1024)

	return nil
}

func (d *Downloader) DownloadAll(repoID string, files []models.FileMetadata) (int, error) {
	downloaded := 0
	for _, file := range files {
		if err := d.DownloadFile(repoID, file.Path); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Failed to download %s: %v\n", file.Path, err)
			continue
		}
		downloaded++
	}
	return downloaded, nil
}
