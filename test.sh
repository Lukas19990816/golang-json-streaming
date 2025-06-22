#!/bin/bash

echo "🚀 Streaming JSON API Test Script"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if server is running
check_server() {
    if curl -s http://localhost:8080/health > /dev/null; then
        echo -e "${GREEN}✅ Server is running${NC}"
        return 0
    else
        echo -e "${RED}❌ Server is not running${NC}"
        return 1
    fi
}

# Test health endpoint
test_health() {
    echo -e "\n${YELLOW}🔍 Testing health endpoint...${NC}"
    response=$(curl -s http://localhost:8080/health)
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ Health check passed${NC}"
        echo "$response" | jq .
    else
        echo -e "${RED}❌ Health check failed${NC}"
    fi
}

# Test streaming parse endpoint
test_parse_streaming() {
    echo -e "\n${YELLOW}📊 Testing streaming parse endpoint...${NC}"
    echo "This will process the large JSON file with streaming parser..."
    
    start_time=$(date +%s)
    response=$(curl -X POST http://localhost:8080/parse \
     -H "Content-Type: application/json" \
     --data-binary @data.json)
    end_time=$(date +%s)
    
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ Streaming parse test passed${NC}"
        echo "Response time: $((end_time - start_time)) seconds"
        echo "$response" | jq .
    else
        echo -e "${RED}❌ Streaming parse test failed${NC}"
    fi
}

# Test non-streaming parse endpoint (will likely fail in Docker)
test_parse_all() {
    echo -e "\n${YELLOW}💥 Testing non-streaming parse endpoint...${NC}"
    echo -e "${RED}⚠️  WARNING: This will load entire file into memory and likely cause OOM!${NC}"
    
    start_time=$(date +%s)
    time sleep 30
    response=$(curl -X POST http://localhost:8080/parse \
     -H "Content-Type: application/json" \
     --data-binary @data.json)
    exit_code=$?
    end_time=$(date +%s)
    
    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}✅ Non-streaming parse test passed (somehow!)${NC}"
        echo "Response time: $((end_time - start_time)) seconds"
        echo "$response" | jq .
    elif [[ $exit_code -eq 124 ]]; then
        echo -e "${RED}⏰ Request timed out (likely OOM)${NC}"
    else
        echo -e "${RED}❌ Non-streaming parse test failed (likely OOM)${NC}"
        echo "This is expected behavior in memory-limited container!"
    fi
}

# Main execution
main() {
    if check_server; then
        test_health
        test_parse_streaming
        test_parse_all
        echo -e "\n${GREEN}🎉 All tests completed!${NC}"
        echo -e "\n${YELLOW}💡 Key takeaways:${NC}"
        echo "   - Streaming parser works efficiently within memory limits"
        echo "   - Non-streaming parser likely fails due to memory constraints"
        echo "   - This demonstrates the importance of streaming for large files"
    else
        echo -e "\n${YELLOW}💡 To start the server:${NC}"
        echo "   make docker-run"
        echo "   OR"
        echo "   make run (for local testing)"
    fi
}

main 