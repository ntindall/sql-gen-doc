#!/usr/bin/env bash

set -euo pipefail
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck source=/dev/null
source "$DIR"/logs.sh

dep=$(which dep)
if [ ! -z "$dep" ]; then
  log "dep is already installed at: $dep (nothing to do)."
  exit 0
fi

log "go get -u github.com/golang/dep/cmd/dep"
go get -u github.com/golang/dep/cmd/dep
