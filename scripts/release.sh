#!/bin/bash
# tags current master as RELEASE_TAG
set -euo pipefail
DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck source=/dev/null
source "$DIR"/logs.sh

: ${RELEASE_TAG:=""}

if [ -z "$RELEASE_TAG" ]; then
  log_error "RELEASE_TAG must be specified"
  exit 1
fi

branch=$(git rev-parse --abbrev-ref HEAD)

if [ $branch != "master" ]; then
  log_error "can only publish tags from the master branch"
  exit 1
fi

git diff-files --quiet ||
  (log_error "Working directory contains unstaged changes"; exit 1)

log "tagging $branch as $RELEASE_TAG"
git tag $RELEASE_TAG $branch
git push --tags

log "successfully published $RELEASE_TAG"
