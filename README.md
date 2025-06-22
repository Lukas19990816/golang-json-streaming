# Streaming JSON API

A memory-efficient Golang API server that can process large JSON files using streaming parser, designed to run in a Docker container with only 50MB memory limit.

## Features

- Streaming JSON parser (memory-friendly)
- Docker deployment with memory constraints
- Real-time memory monitoring
- Health check endpoint
- Progress tracking

## Quick Start

1. **Generate test data (500MB JSON file):**
   ```bash
   make generate
   # OR manually:
   cd cmd/generate && go run main.go
   ```

2. **Build and run with Docker:**
   ```bash
   make docker-run
   # OR manually:
   docker-compose up --build
   ```

3. **Test the API:**
   ```bash
   # Parse with streaming (should work)
   curl http://localhost:8080/parse

   # Parse by loading all (will likely fail in Docker)
   curl http://localhost:8080/parse-all

   # Check health and memory stats
   curl http://localhost:8080/health
   ```

   ```bash
   curl -X POST http://localhost:8080/parse \
     -H "Content-Type: application/json" \
     --data-binary @data.json
   ```

4. **Monitor memory usage:**
   ```bash
   docker stats
   ```

## API Endpoints

- `GET /parse` - Parse the large JSON file using **streaming parser** (memory efficient)
- `GET /parse-all` - Parse the large JSON file by **loading all into memory** (will likely fail with OOM)
- `GET /health` - Health check with memory statistics

## 🔬 Streaming vs Non-streaming Comparison

| Method | Memory Usage | 500MB File | Docker 50MB Limit |
|--------|-------------|------------|-------------------|
| **Streaming** (`/parse`) | ~8-15MB | ✅ Works | ✅ Works |
| **Load All** (`/parse-all`) | ~500MB+ | ❌ OOM | ❌ Container killed |

## Expected Output

**Streaming endpoint (`/parse`):**
```json
{
  "method": "streaming",
  "records_parsed": 1000000,
  "duration_ms": 15234,
  "memory_stats": {
    "alloc": 8,
    "total_alloc": 127,
    "sys": 18,
    "num_gc": 45
  }
}
```

**Non-streaming endpoint (`/parse-all`) in Docker:**
```bash
# Likely results:
curl: (52) Empty reply from server
# OR
curl: (7) Failed to connect to localhost port 8080: Connection refused

# Docker container gets killed due to memory limit exceeded
```

The streaming parser stays well under 50MB while the non-streaming approach gets killed by Docker's memory limit. 