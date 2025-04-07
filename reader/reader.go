//package reader
//
//// reader.go
//// Reads a CSV file line by line and streams URLs through a channel.
//
//import (
//	"bufio"
//	"os"
//	"strings"
//
//	"github.com/rs/zerolog/log"
//)
//
//// ReadCSVFile reads the input CSV file line-by-line (without loading all into memory).
//// It sends valid URLs (with "https://" prefix added if missing) to the provided channel.
//
//func ReadCSVFile(csvPath string, urlChan chan<- string) {
//	// Open the file
//	// Skip the header
//	// For each line: trim, sanitize, and send to channel
//	// Close channel when done
//	defer close(urlChan)
//
//	file, err := os.Open(csvPath)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to open CSV file")
//		return
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	lineNum := 0
//
//	for scanner.Scan() {
//		line := strings.TrimSpace(scanner.Text())
//		lineNum++
//		if lineNum == 1 {
//			continue
//		}
//		if line == "" {
//			continue
//		}
//		if !strings.HasPrefix(line, "http") {
//			line = "https://" + line
//		}
//		urlChan <- line
//	}
//
//	if err := scanner.Err(); err != nil {
//		log.Error().Err(err).Msg("Error scanning CSV file")
//	}
//}

package reader

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// ReadCSV reads URLs line-by-line from a CSV file and sends them through urlChan.
// It expects the first line to be a header. Logs an error if the file is empty or malformed.
func ReadCSVFile(ctx context.Context, filePath string, urlChan chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to open CSV file")
		return err
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	// Read header
	header, err := r.Read()
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to read CSV header, check the file content")
		close(urlChan)
		return err
	}
	if len(header) == 0 || strings.ToLower(header[0]) != "urls" {
		err := errors.New("missing or invalid CSV header")
		log.Error().Err(err).Str("file", filePath).Msg("Invalid CSV format")
		close(urlChan)
		return err
	}

	// Attempt to read first URL to detect empty data
	record, err := r.Read()
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("CSV file is empty or contains no data rows")
		close(urlChan)
		return err // Either EOF or real error — both signal an empty or bad file
	}

	// First valid URL line
	select {
	case <-ctx.Done():
		return nil
	case urlChan <- record[0]:
	}

	// Read remaining lines
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if len(record) == 0 {
			continue
		}
		select {
		case <-ctx.Done():
			return nil
		case urlChan <- record[0]:
		}
	}

	close(urlChan) // Signal no more URLs
	log.Info().Msg("Finished reading CSV file.")
	return nil
}
