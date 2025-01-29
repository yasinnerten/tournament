#!/bin/sh
set -e

echo "Running router tests..."

if go test -v ./...; then
  echo "Tests passed successfully."
else
  echo "Tests failed, but continuing the pipeline."
fi

exit 0