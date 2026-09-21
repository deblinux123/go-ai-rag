#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

mkdir -p dist

echo "==> Building Go AI RAG..."
go build   -trimpath   -ldflags="-s -w"   -o dist/go-ai-rag   ./cmd/app

echo "==> Building Linux desktop package..."
if command -v fyne >/dev/null 2>&1; then
  (
    cd cmd/app
    fyne package --os linux --release
    mv ./*.tar.gz ../../dist/go-ai-rag-linux-amd64.tar.gz
  )
else
  echo "Fyne CLI is not installed; skipping desktop package."
  echo "Install it with:"
  echo "  go install fyne.io/tools/cmd/fyne@latest"
fi

echo
echo "Binary:  $ROOT/dist/go-ai-rag"
echo "Package: $ROOT/dist/go-ai-rag-linux-amd64.tar.gz"
