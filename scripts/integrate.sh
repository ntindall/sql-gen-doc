#!/usr/bin/env bash

set -euo pipefail
DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck source=/dev/null
source "$DIR"/logs.sh

mkdir -p tmp
function finish {
  rm -rf tmp
}
trap finish EXIT

#ensur the db exists
mysql -h mysql -uroot -ppassword -e "CREATE DATABASE IF NOT EXISTS example; GRANT ALL PRIVILEGES ON example.* TO user IDENTIFIED BY 'password'"

log "migrating database down to the beginning"
make migrate-reset
log "bringing database back up"
make migrate-up

log "testing create a new file"
./bin/sql-gen-doc -dsn 'user:password@tcp(mysql:3306)/example' -o tmp/example.md
cp tmp/example.md logs/out1.md
if [ "$(diff --text tmp/example.md fixtures/expected1.md |& tee logs/test1.diff)" ]; then
  log_error "output did not match fixture"
  exit 1
fi

log "testing insert between markdown"
cp fixtures/seed2.md tmp/example2.md
./bin/sql-gen-doc -dsn 'user:password@tcp(mysql:3306)/example' -o tmp/example2.md
cp tmp/example2.md logs/out2.md

if [ "$(diff --text tmp/example2.md fixtures/expected2.md |& tee logs/test2.diff)" ]; then
  log_error "output did not match fixture"
  exit 1
fi

