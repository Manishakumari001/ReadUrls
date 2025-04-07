// writer.go
// Handles writing downloaded content to disk using a single goroutine.

package writer

import (
	"ReadUrls/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

// PersistContent receives byte slices and writes them to randomly named .txt files.
// Only one goroutine should run this function to avoid race conditions.
func PersistContent(contentChan <-chan []byte, outputDir string) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("Failed to create output directory")
	}

	for content := range contentChan {
		fileName := utils.GenerateRandomFileName()
		path := filepath.Join(outputDir, fileName)

		if err := os.WriteFile(path, content, 0644); err != nil {
			log.Error().Err(err).Str("file", path).Msg("Failed to write file")
			continue
		}
		log.Info().Str("file", path).Msg("File saved")
	}
}
