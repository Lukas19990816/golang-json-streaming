package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Entry struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func main() {
	log.Println("Generating large JSON file...")
	
	f, err := os.Create("../../data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Write opening bracket
	f.Write([]byte("[\n"))

	// Generate ~500MB of JSON data
	// Each entry is approximately 500 bytes, so we need ~1M entries
	totalEntries := 1_000_000
	longValue := strings.Repeat("This is a long string to make the JSON file larger. ", 10)

	for i := 0; i < totalEntries; i++ {
		entry := Entry{
			ID:    i,
			Value: fmt.Sprintf("%s Entry number %d", longValue, i),
		}
		
		entryJSON, _ := json.Marshal(entry)
		f.Write(entryJSON)
		
		if i != totalEntries-1 {
			f.Write([]byte(",\n"))
		} else {
			f.Write([]byte("\n"))
		}
		
		// Progress indication
		if i%100000 == 0 {
			log.Printf("Generated %d/%d entries", i, totalEntries)
		}
	}

	// Write closing bracket
	f.Write([]byte("]"))

	// Check file size
	info, _ := f.Stat()
	sizeMB := float64(info.Size()) / 1024 / 1024
	log.Printf("Generated data.json with size: %.2f MB", sizeMB)
} 