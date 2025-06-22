package main

import (
	"encoding/json"
	"log"
	"os"
)

// 每個 entry 的結構（value 是 slice of strings）
type Entry struct {
	ID    int    `json:"-"`
	Value []Node `json:"value"`
}

type Node struct {
	NodeId   string `json:"nodeId"`
	NodeName string `json:"nodeName"`
}

func main() {
	log.Println("Generating large JSON object...")

	f, err := os.Create("../../data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 寫入 JSON 物件的開頭
	f.WriteString("{\n")
	f.WriteString("\"id\": 1,\n")
	f.WriteString("\"value\":[\n\t")

	// longValue := strings.Repeat("This is a long string to make the JSON file larger. ", 10)

	for i := 0; ; i++ {
		node := Node{
			NodeId:   string(int(i)),
			NodeName: "node" + string(i),
		}

		// 使用 json.Encoder 流式寫入 value 部分（slice of strings）
		encoder := json.NewEncoder(f)
		if err := encoder.Encode(node); err != nil {
			log.Fatalf("Failed to encode entry: %v", err)
		}

		// 判斷是否為最後一個項目，避免多個 comma
		_, err := f.Seek(0, 2) // 獲取當前文件位置
		if err != nil {
			log.Fatal(err)
		}
		posAfterEncoding, err := f.Seek(0, 1)
		if err != nil {
			log.Fatal(err)
		}

		// 如果不是最後一個 entry，寫入逗號和換行
		if posAfterEncoding < 1*1024*1024 { // 假設 500MB 為限界點
			f.WriteString(",\n\t")
		} else {
			break
		}
	}

	// 寫入 JSON 物件的結尾
	f.WriteString("]\n")
	f.WriteString("}")
}
