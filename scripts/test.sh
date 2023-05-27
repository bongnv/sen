#!/bin/bash

set -e

find . -name "go.mod" | while read -r fname; do
  module=$(dirname "$fname")
  echo "Testing ${module}"
  pushd "$module" > /dev/null
  go test -race --coverprofile=coverage.coverprofile ./...
  popd > /dev/null
done
