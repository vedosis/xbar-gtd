#!/bin/bash
set -e

echo "🔍 Pre-commit checks..."

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -n "$STAGED_GO_FILES" ]; then
    echo "📝 Formatting Go files..."

    # Format Go files
    go fmt ./...

    # Re-stage formatted files
    echo "$STAGED_GO_FILES" | xargs git add

    echo "🏗️  Updating BUILD files..."

    # Update BUILD files with gazelle
    bazel run //:gazelle

    # Stage updated BUILD files
    git add **/BUILD.bazel BUILD.bazel
fi

# Check if go.mod or go.sum changed
GO_MOD_CHANGED=$(git diff --cached --name-only | grep -E '^go\.(mod|sum)$' || true)

if [ -n "$GO_MOD_CHANGED" ]; then
    echo "📦 Tidying Go modules..."
    go mod tidy

    # Re-stage go.mod and go.sum
    git add go.mod go.sum

    echo "🔄 Updating Bazel dependencies..."
    bazel run //:gazelle-update-repos
fi

echo "✅ Pre-commit checks passed!"