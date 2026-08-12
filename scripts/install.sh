#!/usr/bin/env bash
set -euo pipefail

# 第1引数があれば優先し、無ければ Git タグを自動取得。タグが取得できなければ "dev" に設定します
VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
  VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
fi

COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
  COMMIT="${COMMIT}-dirty"
fi
DATE=$(date --iso-8601=seconds 2>/dev/null || date +'%Y-%m-%dT%H:%M:%S%z')

PKG="github.com/MasayukiFukada/PjDoc/internal/features/show_version"
LDFLAGS="-X ${PKG}.Version=${VERSION} -X ${PKG}.Commit=${COMMIT} -X ${PKG}.Date=${DATE}"

echo "Installing PjDoc with version: ${VERSION}, commit: ${COMMIT}, date: ${DATE}"
go install -ldflags "${LDFLAGS}" ./cmd/pjdoc
