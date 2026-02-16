.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: lint-fast
lint-fast:
	golangci-lint run --fast
