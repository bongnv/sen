#!/bin/bash

set -e

find . -name "go.mod" | while read -r fname; do
  module=$(dirname "$fname")
  echo "Linting ${module}"
  pushd "$module" > /dev/null
  golangci-lint run $@
  popd > /dev/null
done
