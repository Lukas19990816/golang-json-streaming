package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Entry struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
}

func getMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemoryStats{
		Alloc:      m.Alloc / 1024 / 1024,      // MB
		TotalAlloc: m.TotalAlloc / 1024 / 1024, // MB
		Sys:        m.Sys / 1024 / 1024,        // MB
		NumGC:      m.NumGC,
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	file, err := os.Open("data.json")
	if err != nil {
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// 開始解析 JSON 陣列
	t, err := decoder.Token()
	if err != nil || t != json.Delim('[') {
		http.Error(w, "invalid JSON array", http.StatusBadRequest)
		return
	}

	count := 0
	for decoder.More() {
		var entry Entry
		if err := decoder.Decode(&entry); err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}
		count++

		// Force garbage collection every 10000 records
		if count%10000 == 0 {
			runtime.GC()
		}
	}

	// Read the closing delimiter
	t, err = decoder.Token()
	if err != nil || t != json.Delim(']') {
		http.Error(w, "invalid JSON array end", http.StatusBadRequest)
		return
	}

	duration := time.Since(start)
	memStats := getMemoryStats()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"method":         "streaming",
		"records_parsed": count,
		"duration_ms":    duration.Milliseconds(),
		"memory_stats":   memStats,
	}

	json.NewEncoder(w).Encode(response)
}

// Non-streaming handler - loads entire JSON into memory at once
func parseAllHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	log.Println("⚠️  WARNING: Using non-streaming parser - will load entire file into memory!")

	// Read entire file into memory
	data, err := os.ReadFile("data.json")
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	log.Printf("📁 File size: %.2f MB", float64(len(data))/1024/1024)
	memStatsAfterRead := getMemoryStats()
	log.Printf("📊 Memory after file read: %+v", memStatsAfterRead)

	// Parse entire JSON array into memory
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		http.Error(w, "failed to parse JSON", http.StatusInternalServerError)
		return
	}

	count := len(entries)
	duration := time.Since(start)

	// Force GC to see real memory usage
	runtime.GC()
	runtime.GC() // Call twice for better measurement

	memStats := getMemoryStats()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"method":         "load-all",
		"records_parsed": count,
		"duration_ms":    duration.Milliseconds(),
		"memory_stats":   memStats,
		"file_size_mb":   float64(len(data)) / 1024 / 1024,
	}

	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	memStats := getMemoryStats()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":       "ok",
		"memory_stats": memStats,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/parse", streamHandler)
	http.HandleFunc("/parse-all", parseAllHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Server starting on :8080")
	log.Println("Endpoints:")
	log.Println("  GET /parse     - Streaming JSON parser (memory efficient)")
	log.Println("  GET /parse-all - Load all JSON into memory (will likely OOM)")
	log.Println("  GET /health    - Health check with memory stats")
	log.Printf("Initial memory stats: %+v", getMemoryStats())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
