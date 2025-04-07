package reader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestReadCSVFile_Empty(t *testing.T) {
	// Create empty test file
	filePath := filepath.Join(os.TempDir(), "empty.csv")
	os.WriteFile(filePath, []byte("Urls\n"), 0644)
	ctx := context.Background()

	urlChan := make(chan string)
	go ReadCSVFile(ctx, filePath, urlChan)

	count := 0
	for range urlChan {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 URLs, got %d", count)
	}
}

func TestReadCSVFile_Valid(t *testing.T) {
	ctx := context.Background()
	data := "Urls\nwww.example.com\nhttps://another.com\n"
	filePath := filepath.Join(os.TempDir(), "valid.csv")
	os.WriteFile(filePath, []byte(data), 0644)

	urlChan := make(chan string)
	go ReadCSVFile(ctx, filePath, urlChan)

	var urls []string
	for url := range urlChan {
		urls = append(urls, url)
	}

	expected := []string{"https://www.example.com", "https://another.com"}
	for i, want := range expected {
		if urls[i] != want {
			t.Errorf("Expected %s, got %s", want, urls[i])
		}
	}
}
