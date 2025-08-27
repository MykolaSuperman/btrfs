.PHONY: build clean docker-build docker-run test help

help:
	@echo "Available targets:"
	@echo "  build        - Build the Go package"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  test         - Run tests (requires BTRFS)"

build:
	go build

build-cli:
	go build -o btrfs-cli cmd/main.go

clean:
	go clean
	rm -f btrfs-cli

docker-build:
	docker build -t btrfs-go .

docker-run:
	docker run --privileged -it btrfs-go

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

help:
	@echo "BTRFS Go Package - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build        - Build the package"
	@echo "  make build-cli    - Build the CLI tool"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make test         - Run tests"
	@echo "  make deps         - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make help         - Show this help"
