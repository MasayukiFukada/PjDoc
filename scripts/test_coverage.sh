#!/usr/bin/env bash

# エラーが発生した場合は即座に停止します
set -e

echo "🧪 Running unit tests and generating coverage profile..."
go test -coverprofile=coverage.out ./...

echo "📊 Test suite passed! Opening interactive HTML coverage report in your browser..."
go tool cover -html=coverage.out
