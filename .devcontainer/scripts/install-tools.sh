#!/bin/sh

set -e  # exit immediately on error
set -u  # error on undefined variables

echo "Installing Go development tools..."

# Retries a `go install` a few times to ride out transient network/sumdb
# errors, and skips it entirely if the binary is already on PATH.
# Package spec is read from the GOFUMPT_PKG-style env var below rather
# than passed as a second CLI word, so a stray line-wrap can't drop it.
install_tool() {
  bin_name="$1"

  if command -v "$bin_name" >/dev/null 2>&1; then
    echo "  $bin_name already installed, skipping."
    return 0
  fi

  i=1
  while [ "$i" -le 3 ]; do
    if go install "$PKG"; then
      return 0
    fi
    echo "  install of $PKG failed (attempt $i/3), retrying in 5s..."
    sleep 5
    i=$((i + 1))
  done

  echo "  failed to install $PKG after 3 attempts"
  return 1
}

echo "Installing gofumpt..."
PKG="mvdan.cc/gofumpt@v0.11.0"
install_tool gofumpt

echo "Installing gci..."
PKG="github.com/daixiang0/gci@v0.14.0"
install_tool gci

echo "Installing golangci-lint..."
PKG="github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"
install_tool golangci-lint

echo "Installing lefthook..."
PKG="github.com/evilmartians/lefthook/v2@v2.1.11"
install_tool lefthook

echo "All Go dev tools installed!"