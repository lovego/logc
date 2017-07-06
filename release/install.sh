#!/bin/bash

set -e
if ! which logc; then
  sudo wget -O /usr/local/bin/logc 'https://github.com/lovego/logc/releases/download/170706/logc'
  sudo chmod +x /usr/local/bin/logc
fi
