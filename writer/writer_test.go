package writer

import (
	"os"
	"testing"
)

func TestPersistContent_WritesFile(t *testing.T) {
	tempDir := t.TempDir()
	contentChan := make(chan []byte, 1)
	contentChan <- []byte("sample text")
	close(contentChan)

	PersistContent(contentChan, tempDir)

	// Check that one .txt file exists
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read tempDir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, found %d", len(files))
	}
}
