//// downloader.go
//// Manages concurrent URL downloads using up to 50 goroutines.
//
//package downloader
//
//import (
//	"context"
//	"github.com/rs/zerolog/log"
//	"io"
//	"net/http"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//// DownloadManager reads URLs from urlChan, downloads their content
//// concurrently using up to 50 goroutines, and sends the result to contentChan.
//// It tracks success/failure counts and total download time.
//func DownloadManager(
//	ctx context.Context,
//	urlChan <-chan string,
//	contentChan chan<- []byte,
//	successCount, failureCount, totalDuration *int64,
//) {
//	var wg sync.WaitGroup
//	sem := make(chan struct{}, 50) // Semaphore to limit concurrency
//
//	// Use a separate goroutine to consume from urlChan and launch downloaders
//	go func() {
//		for url := range urlChan {
//			sem <- struct{}{} // Acquire semaphore
//			wg.Add(1)
//
//			go func(url string) {
//				defer func() {
//					<-sem // Release semaphore
//					wg.Done()
//				}()
//
//				start := time.Now()
//				data, err := downloadURL(ctx, url)
//				if err != nil {
//					log.Error().Err(err).Str("url", url).Msg("Download failed")
//					atomic.AddInt64(failureCount, 1)
//					return
//				}
//
//				atomic.AddInt64(successCount, 1)
//				atomic.AddInt64(totalDuration, int64(time.Since(start)))
//				contentChan <- data
//			}(url)
//		}
//		// All URLs have been scheduled, wait for downloads
//		wg.Wait()
//		close(contentChan)
//	}()
//
//	// Block until context is cancelled
//	<-ctx.Done()
//}
//
//// / downloadURL performs an HTTP GET request and returns the body as bytes.
//// Returns error on failure.
//func downloadURL(ctx context.Context, url string) ([]byte, error) {
//	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil || resp.StatusCode >= 400 {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	return io.ReadAll(resp.Body)
//}

package downloader

import (
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// DownloadManager reads URLs from urlChan, downloads their content
// concurrently using up to 50 goroutines, and sends the result to contentChan.
// It tracks success/failure counts and total download time.
func DownloadManager(
	ctx context.Context,
	urlChan <-chan string,
	contentChan chan<- []byte,
	successCount, failureCount, totalDuration *int64,
	done chan<- struct{}, // new channel to notify completion
) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // Semaphore to limit concurrency

	go func() {
		for url := range urlChan {
			sem <- struct{}{} // Acquire semaphore
			wg.Add(1)

			go func(url string) {
				defer func() {
					<-sem // Release semaphore
					wg.Done()
				}()

				start := time.Now()
				data, err := downloadURL(ctx, url)
				if err != nil {
					log.Error().Err(err).Str("url", url).Msg("Download failed")
					atomic.AddInt64(failureCount, 1)
					return
				}

				atomic.AddInt64(successCount, 1)
				atomic.AddInt64(totalDuration, int64(time.Since(start)))
				contentChan <- data
			}(url)
		}
		// All URLs have been scheduled, wait for downloads
		wg.Wait()
		close(contentChan) // signal writer to finish
		done <- struct{}{} // notify main that downloading is complete
	}()
}

// downloadURL performs an HTTP GET request and returns the body as bytes.
// Returns error on failure.
func downloadURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
