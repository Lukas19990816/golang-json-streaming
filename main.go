package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

type Entry struct {
	ID    int    `json:"id"`
	Value []Node `json:"value"`
	Pod   []Pod  `json:"pod"`
}

type Node struct {
	NodeId   string `json:"nodeId"`
	NodeName string `json:"nodeName"`
}

type Pod struct {
	PodId   string `json:"podId"`
	PodName string `json:"podName"`
}

type MemoryStats struct {
	Alloc      float64 `json:"alloc"`
	TotalAlloc float64 `json:"total_alloc"`
	Sys        float64 `json:"sys"`
	NumGC      uint32  `json:"num_gc"`
}

func getMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemoryStats{
		Alloc:      float64(m.Alloc) / (1024 * 1024),      // MB as float64
		TotalAlloc: float64(m.TotalAlloc) / (1024 * 1024), // MB as float64
		Sys:        float64(m.Sys) / (1024 * 1024),        // MB as float64
		NumGC:      m.NumGC,
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("receive streaming request")

	// 使用 r.Body 來流式讀取 JSON 資料，避免一次性載入記憶體
	decoder := json.NewDecoder(r.Body)

	t, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	count := 0

	for decoder.More() {
		var Node Node
		if err := decoder.Decode(&Node); err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}

		log.Println(Node)

		count++

		// 強制 GC 每處理 10k 筆資料，幫助釋放記憶體
		if count%10000 == 0 {
			runtime.GC()
		}

		stats := getMemoryStats()
		log.Printf("Processed %d records. Current memory usage: %.2fMB", count, stats.Alloc)
	}

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	for decoder.More() {
		var Pod Pod
		if err := decoder.Decode(&Pod); err != nil {
			http.Error(w, "decode error", http.StatusInternalServerError)
			return
		}

		log.Println(Pod)

		count++

		// 強制 GC 每處理 10k 筆資料，幫助釋放記憶體
		if count%10000 == 0 {
			runtime.GC()
		}

		stats := getMemoryStats()
		log.Printf("Processed %d records. Current memory usage: %.2fMB", count, stats.Alloc)
	}

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%T: %v\n", t, t)

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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	log.Printf("📁 File size: %.2f MB", float64(len(data))/1024/1024)
	memStatsAfterRead := getMemoryStats()
	log.Printf("📊 Memory after file read: %+v", memStatsAfterRead)

	// Parse entire JSON array into memory
	var entries Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		http.Error(w, "failed to parse JSON", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)

	// Force GC to see real memory usage
	runtime.GC()
	runtime.GC() // Call twice for better measurement

	memStats := getMemoryStats()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"method":       "load-all",
		"duration_ms":  duration.Milliseconds(),
		"memory_stats": memStats,
		"file_size_mb": float64(len(data)) / 1024 / 1024,
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
