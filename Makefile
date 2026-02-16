.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: lint-fast
lint-fast:
	golangci-lint run --fast

.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/hexlet-path-size ./cmd/hexlet-path-size
	@echo "Build complete: bin/hexlet-path-size"

.PHONY: run
run:
	@go run ./cmd/hexlet-path-size

.PHONY: clean
clean:
	@rm -rf bin
	@echo "Cleaned up bin directory"
