#!/bin/bash

RED="\033[31m"
GREEN="\033[32m"
RESET="\033[0m"

log() {
  echo -e "$GREEN[info] $1$RESET"
}

log_error() {
  echo -e "$RED[error] $1$RESET" 1>&2;
}
