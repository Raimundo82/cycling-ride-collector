#!/bin/sh

set -e  # exit immediately on error
set -u  # error on undefined variables

echo "🔧 Installing Go development tools..."

# Strict gofmt replacement
echo "Installing gofumpt..."
go install mvdan.cc/gofumpt@latest

# Import grouping formatter
echo "Installing gci..."
go install github.com/daixiang0/gci@latest

# Linter
echo "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2
# Hook manager
echo "Installing lefthook..."
go install github.com/evilmartians/lefthook/v2@v2.1.11

echo "✅ All Go dev tools installed!"
