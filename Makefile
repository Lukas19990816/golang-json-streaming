.PHONY: generate build run docker-build docker-run clean test

# Generate test data
generate:
	@echo "Generating 500MB test JSON file..."
	cd cmd/generate && go run main.go

# Build the main application
build:
	go build -o bin/app main.go

# Run locally
run: build
	./bin/app

# Build docker image
docker-build:
	podman-compose build

# Run with docker
docker-run: generate docker-build
	docker-compose up

# Test the API
test:
	./test.sh

# Compare streaming vs non-streaming (for demo purposes)
compare:
	@echo "🔬 Comparing Streaming vs Non-streaming parsers"
	@echo "=============================================="
	@echo ""
	@echo "1️⃣ Testing streaming parser..."
	@curl -s http://localhost:8080/parse | jq .
	@echo ""
	@echo "2️⃣ Testing non-streaming parser (will likely fail in Docker)..."
	@timeout 15 curl -s http://localhost:8080/parse-all | jq . || echo "❌ Failed as expected - OOM!"

# Clean up
clean:
	rm -f data.json
	rm -rf bin/
	docker-compose down --rmi all

# Full setup and run
all: generate docker-run 