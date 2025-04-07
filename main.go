// main.go
// Entry point for the URLFetcher CLI application.
// This file initializes the pipeline, handles command-line input,
// and sets up graceful shutdown with cancellation.

package main

import (
	"ReadUrls/downloader"
	"ReadUrls/reader"
	"ReadUrls/writer"
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Setup structured logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Ensure a CSV file path is provided
	if len(os.Args) < 2 {
		log.Fatal().Msg("Please provide path to CSV file as argument")
	}
	csvFilePath := os.Args[1]

	log.Info().Msg("Starting URL Fetcher")

	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for system interrupt signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channels for pipeline stages
	urlChan := make(chan string)
	contentChan := make(chan []byte)
	done := make(chan struct{})

	// Metrics counters
	var successCount int64
	var failureCount int64
	var totalDuration int64

	// Start CSV reader goroutine (Stage 1)
	go reader.ReadCSVFile(ctx, csvFilePath, urlChan)

	// Start downloader manager goroutine (Stage 2)
	go downloader.DownloadManager(ctx, urlChan, contentChan, &successCount, &failureCount, &totalDuration, done)

	// Start writer goroutine (Stage 3)
	go writer.PersistContent(contentChan, "downloads")

	// Wait for signal to shutdown
	select {
	case <-sigChan:
		log.Warn().Msg("Shutdown signal received. Waiting 5s for downloads to complete...")
		time.Sleep(5 * time.Second)
		cancel()
	case <-ctx.Done():
	// Normal context cancellation
	case <-done:
		//Program terminates after reading & writing all Urls
		time.Sleep(2 * time.Second)
		cancel()

	}

	// Log metrics
	log.Info().
		Int64("total_urls", successCount+failureCount).
		Int64("successful", successCount).
		Int64("failed", failureCount).
		Dur("avg_download_duration", time.Duration(totalDuration/max(successCount, 1))).
		Msg("Processing complete")
}

// max returns the greater of two int64 values
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

