# --- Formatting (manual autofix) ---
format:
	gofumpt -w .
	gci write --skip-generated .

# --- Formatting checks (no autofix, used in hooks + CI) ---
check-format:
	gofumpt -l .
	gci list --skip-generated .

# --- Linting ---
lint:
	golangci-lint run

# --- Tests ---
test:
	go test ./... -v

test-integration:
	go test -tags=integration -v ./internal/tests/integration/...

test-all: test test-integration

# --- Combined checks for pre-commit ---
precommit: check-format lint test

# --- Combined checks for pre-push ---
prepush: check-format lint test
