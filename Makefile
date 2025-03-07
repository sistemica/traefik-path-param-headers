.PHONY: build test lint dev clean

# Build the plugin
build:
	go build -v ./...

# Run tests
test:
	go test -v ./...

# Run linting
lint:
	golangci-lint run

# Start development environment
dev:
	cd dev && docker-compose up

# Clean build artifacts
clean:
	rm -rf dist/
	
# Create a new release
release:
	@echo "Creating release..."
	@read -p "Enter version (e.g. v1.0.0): " VERSION; \
	git tag $$VERSION && \
	git push origin $$VERSION

# Initialize a new development environment
init:
	@echo "Initializing development environment..."
	mkdir -p .assets
	touch .assets/icon.png
	touch .assets/banner.png
	mkdir -p .github/workflows