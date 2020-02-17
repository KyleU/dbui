#!/bin/bash

## Runs code statistics, checks for outdated dependencies, then runs linters

set -e
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
project_dir=${dir}/..
cd $project_dir

echo "=== linting ==="
golangci-lint run \
  -E deadcode \
  -E errcheck \
  -E gosimple \
  -E govet \
  -E ineffassign \
  -E structcheck \
  -E typecheck \
  -E unused \
  -E varcheck \
  -E bodyclose \
  -E depguard \
  -E dogsled \
  -E dupl \
  -E funlen \
  -D gochecknoglobals \
  -E gochecknoinits \
  -E gocognit \
  -E goconst \
  -E gocritic \
  -E gocyclo \
  -E godox \
  -E gofmt \
  -E goimports \
  -E golint \
  -E gomnd \
  -E goprintffuncname \
  -D gosec \
  -E interfacer \
  -E lll \
  -E maligned \
  -E misspell \
  -E nakedret \
  -E prealloc \
  -E rowserrcheck \
  -E scopelint \
  -E stylecheck \
  -E unconvert \
  -E unparam \
  -E whitespace \
  -D wsl \
./...

