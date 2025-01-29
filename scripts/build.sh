#!/bin/sh
set -e # exit immediately if a command exits with a non-zero status

set -x # Print commands and their arguments as they are executed

echo "Building the app..."

mkdir -p ./bin
mkdir -p ./cmd/app

# Initialize the Go module if not already initialized
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init tournament-app
fi

go mod tidy
go mod vendor

go build -mod=vendor -o ./bin/app ./cmd/app/main.go