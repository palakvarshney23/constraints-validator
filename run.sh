#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Install from https://go.dev/dl/ then run again."
    exit 1
fi
go run .
