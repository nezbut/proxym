# =================================================================================================
# Linters
# =================================================================================================

lint:
	sh scripts/lint.sh

lint-fix:
	@echo "Fixing lint issues..."
	golangci-lint run --fix ./...

# =================================================================================================
# Tests
# =================================================================================================

test:
	@echo "Running tests..."
	go test -v ./...
