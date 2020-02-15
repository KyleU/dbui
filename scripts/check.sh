#!/bin/bash

## Runs code statistics, checks for outdated dependencies, then runs linters

set -e
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
project_dir=${dir}/..
cd $project_dir

echo "=== todo ==="
