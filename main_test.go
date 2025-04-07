package main

import (
	"ReadUrls/downloader"
	"ReadUrls/reader"
	"ReadUrls/writer"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPipelineIntegration(t *testing.T) {
	// Setup fake CSV input
	content := "Urls\nhttps://example.com\n"
	filePath := filepath.Join(os.TempDir(), "pipeline.csv")
	os.WriteFile(filePath, []byte(content), 0644)

	ctx := context.Background()
	urlChan := make(chan string)
	contentChan := make(chan []byte)
	done := make(chan struct{})

	var successCount int64
	var failureCount int64
	var totalDuration int64

	go reader.ReadCSVFile(ctx, filePath, urlChan)
	go downloader.DownloadManager(ctx, urlChan, contentChan, &successCount, &failureCount, &totalDuration, done)

	tempDir := t.TempDir()
	writer.PersistContent(contentChan, tempDir)

	if successCount == 0 {
		t.Error("Expected at least 1 successful download")
	}
}
