#!/bin/sh
 
set -e  # exit immediately on error
set -u  # error on undefined variables
 
echo "🔧 Installing Go development tools..."
 
# Retries a `go install pkg@version` a few times to ride out transient
# network/sumdb errors (e.g. sum.golang.org HTTP/2 stream resets), and
# skips the call entirely if the binary is already on PATH.
install_tool() {
  bin_name="$1"
  pkg_at_version="$2"
 
  if command -v "$bin_name" >/dev/null 2>&1; then
    echo "  $bin_name already installed, skipping."
    return 0
  fi
 
  i=1
  while [ "$i" -le 3 ]; do
    if go install "$pkg_at_version"; then
      return 0
    fi
    echo "  ⚠️  install of $pkg_at_version failed (attempt $i/3), retrying in 5s..."
    sleep 5
    i=$((i + 1))
  done
 
  echo "  ❌ failed to install $pkg_at_version after 3 attempts"
  return 1
}

# Strict gofmt replacement
echo "Installing gofumpt..."
install_tool install mvdan.cc/gofumpt@v0.11.0

# Import grouping formatter
echo "Installing gci..."
install_tool install github.com/daixiang0/gci@v0.14.0

# Linter
echo "Installing golangci-lint..."
install_tool install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2

# Hook manager
echo "Installing lefthook..."

install_tool install github.com/evilmartians/lefthook/v2@v2.1.11

echo "✅ All Go dev tools installed!"
