# live-webrtc-go Makefile

.PHONY: help build fmt-check fmt lint vet security test-unit test-integration test-security test-e2e test-performance test test-all coverage clean docker-build docker-run install-tools ci dev load-test benchmark

BINARY_NAME=live-webrtc-go
DOCKER_IMAGE=live-webrtc-go:latest
GO_VERSION=1.22

help:
	@echo "live-webrtc-go - build, test, and verification commands"
	@echo ""
	@echo "Development Commands:"
	@echo "  make build             - Build the application"
	@echo "  make fmt               - Format code with gofmt -w"
	@echo "  make fmt-check         - Check formatting without modifying files"
	@echo "  make lint              - Run vet and golangci-lint"
	@echo "  make security          - Run gosec scan"
	@echo "  make test              - Run default verification tests"
	@echo "  make test-all          - Run all supported test suites"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make test-security     - Run security tests only"
	@echo "  make test-e2e          - Run end-to-end tests only"
	@echo "  make test-performance  - Run performance tests only"
	@echo "  make coverage          - Generate coverage report"
	@echo "  make clean             - Clean build artifacts"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-run        - Run Docker container"
	@echo ""
	@echo "CI/CD Commands:"
	@echo "  make ci                - Run the local CI verification pipeline"
	@echo "  make install-tools     - Install development tools"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Build completed: bin/$(BINARY_NAME)"

install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/AlekSi/gocov-xml@latest
	@echo "Development tools installed"

fmt:
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "Code formatting completed"

fmt-check:
	@echo "Checking formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'make fmt'."; \
		gofmt -s -l .; \
		exit 1; \
	fi
	@echo "Formatting check completed"

vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Go vet completed"

lint: fmt-check vet
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@echo "Linting completed"

security:
	@echo "Running security analysis..."
	@mkdir -p reports
	@gosec -fmt=json -out=reports/security-report.json ./...
	@echo "Security analysis completed - report: reports/security-report.json"

test-unit:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage-unit.out ./internal/...
	@echo "Unit tests completed"

test-integration:
	@echo "Running integration tests..."
	@go test -v -race -tags=integration ./test/integration/...
	@echo "Integration tests completed"

test-security:
	@echo "Running security tests..."
	@go test -v -tags=security ./test/security/...
	@echo "Security tests completed"

test-e2e:
	@echo "Running end-to-end tests..."
	@go test -v -tags=e2e -timeout=10m ./test/e2e/...
	@echo "End-to-end tests completed"

test-performance:
	@echo "Running performance tests..."
	@go test -v -tags=performance -bench=. -benchmem ./test/performance/...
	@echo "Performance tests completed"

test: test-unit test-integration test-security
	@echo "Default test suite completed"

test-all: test-unit test-integration test-security test-e2e test-performance
	@echo "All supported tests completed"

coverage:
	@echo "Generating coverage report..."
	@mkdir -p reports
	@go test -v -race -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o reports/coverage.html
	@gocov convert coverage.out | gocov-xml > reports/coverage.xml
	@echo "Coverage report generated:"
	@echo "  HTML: reports/coverage.html"
	@echo "  XML:  reports/coverage.xml"
	@go tool cover -func=coverage.out | tail -1

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf reports/
	@rm -f coverage*.out
	@rm -f *.test
	@echo "Cleanup completed"

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run:
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 \
		-e RECORD_ENABLED=1 \
		-e RECORD_DIR=/records \
		-v "$$(pwd)/records:/records" \
		$(DOCKER_IMAGE)

ci: clean lint security test coverage
	@echo "Local CI verification completed successfully"

dev:
	@echo "Starting development server..."
	@./scripts/start.sh

load-test:
	@echo "Running load tests..."
	@go run ./test/load -duration=30s -concurrent=10
	@echo "Load tests completed"

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "Benchmarks completed"
