package utils

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenerateRandomFileName(t *testing.T) {
	filename := GenerateRandomFileName()

	// Check that the filename ends with .txt
	if !strings.HasSuffix(filename, ".txt") {
		t.Errorf("Filename does not have .txt extension: %s", filename)
	}

	// Check the full format: 8 alphanumeric characters + .txt
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]{8}\.txt$`)
	if !pattern.MatchString(filename) {
		t.Errorf("Filename format invalid: %s", filename)
	}

	// Generate multiple filenames and check for uniqueness
	fileSet := make(map[string]bool)
	for i := 0; i < 100; i++ {
		name := GenerateRandomFileName()
		if fileSet[name] {
			t.Errorf("Duplicate filename generated: %s", name)
		}
		fileSet[name] = true
	}
}
