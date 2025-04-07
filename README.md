## ReadUrls Application Details

This application is for reading the Urls . It is a high-performance command-line application written in Go that reads a list of URLs from a CSV file, downloads their content concurrently, and persists the content as `.txt` files on disk.
It is designed to be memory-efficient and robust, using a multi-stage pipeline:
- Stage 1: Reads the CSV line-by-line (not fully into memory)
- Stage 2: Downloads content using up to 50 goroutines
- Stage 3: Writes files using a single writer goroutine

This ensures high concurrency without resource exhaustion, and avoids file-write race conditions.

---

### Features

-  Stream-based CSV reader
-  Dynamic goroutine creation (max 50)
-  Single writer to avoid disk write conflicts
-  Graceful shutdown with 5-second timeout
-  Success/failure metrics with timing
-  Unit + integration tests


###  Building

Make sure you have **Go 1.23+** installed.

1. Clone the repository:

```bash
git clone https://github.com/Manishakumari001/ReadUrls.git
cd ReadUrls

go build -o ReadUrls main.go
```

### Run Application
```bash
go run main.go "csv-file-path"

OR

ReadUrls "csv-file-path"
```
