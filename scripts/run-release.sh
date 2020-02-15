#!/bin/bash

## Builds the project in release mode and runs it

set -e
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
project_dir=${dir}/..
cd $project_dir

bin/build-css
# bin/build-client
make build-release
RUST_BACKTRACE=1 build/release/dbui
