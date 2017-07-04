#!/bin/bash

set -e
if ! which logc; then
  sudo wget -O /usr/local/bin/logc "http://github.com/lovego/logc/raw/master/release/logc"
  sudo chmod +x /usr/local/bin/logc
fi

test -d logc || mkdir logc
nohup logc "$@" >>logc/console.log 2>>&1 &
