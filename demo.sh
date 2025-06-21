#!/bin/bash

echo "🚀 Streaming JSON API Demo"
echo "=========================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Step 1: Generate 500MB test JSON file${NC}"
echo "This will create a large JSON file with 1M records..."
make generate

echo -e "\n${BLUE}Step 2: Test locally first${NC}"
echo "Building and running locally..."
make build

echo -e "\n${YELLOW}Starting server in background...${NC}"
./bin/app &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo -e "\n${BLUE}Step 3: Test the API${NC}"
./test.sh

# Kill local server
kill $SERVER_PID 2>/dev/null

echo -e "\n${BLUE}Step 4: Docker deployment with memory limit${NC}"
echo "Building and running with Docker (50MB memory limit)..."
docker-compose up --build -d

echo -e "\n${YELLOW}Waiting for Docker container to start...${NC}"
sleep 5

echo -e "\n${BLUE}Step 5: Test Docker deployment${NC}"
./test.sh

echo -e "\n${BLUE}Step 6: Monitor memory usage${NC}"
echo "Docker container memory stats:"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"

echo -e "\n${GREEN}🎉 Demo completed!${NC}"
echo -e "${YELLOW}To clean up:${NC} make clean" 