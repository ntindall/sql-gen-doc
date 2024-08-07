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

# https://github.com/docker/compose/issues/4369
log "sleeping... docker-compose run does not respect health checks"
sleep 10
#ensure the db exists
mysql -h mysql -uroot -ppassword -e "CREATE DATABASE IF NOT EXISTS example; GRANT ALL PRIVILEGES ON example.* TO user IDENTIFIED BY 'password'"


log "migrating database down to the beginning"
make migrate-reset || echo "ERROR" # allow fail through because the command will return an error if no tables exist

log "bringing database back up"
make migrate-up

had_err=0

# test case 1
log "testing create a new file"
./bin/sql-gen-doc -dsn 'user:password@tcp(mysql:3306)/example' -o tmp/example1.md
cp tmp/example1.md logs/out1.md
if [ "$(diff --text tmp/example1.md fixtures/expected1.md |& tee logs/test1.diff)" ]; then
  cat tmp/example1.md
  log_error "output did not match fixture -- see logs/test1.diff"
  had_err=1
fi

# test case 2
log "testing inserting between markdown"
cp fixtures/seed2.md tmp/example2.md
./bin/sql-gen-doc -dsn 'user:password@tcp(mysql:3306)/example' -o tmp/example2.md
cp tmp/example2.md logs/out2.md

if [ "$(diff --text tmp/example2.md fixtures/expected2.md |& tee logs/test2.diff)" ]; then
  cat tmp/example2.md
  log_error "output did not match fixture -- see logs/test2.diff"
  had_err=1
fi

if [ $had_err ]; then
  exit 0
fi

log "OK"

