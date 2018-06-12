# Logc
A log collector(like logstash, fluentd) written by golang.
Now only support read from files and output to elasticSearch.

[![Build Status](https://travis-ci.org/lovego/logc.svg?branch=master)](https://travis-ci.org/lovego/logc)
[![Coverage Status](https://coveralls.io/repos/github/lovego/logc/badge.svg?branch=master)](https://coveralls.io/github/lovego/logc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/logc)](https://goreportcard.com/report/github.com/lovego/logc)
[![GoDoc](https://godoc.org/github.com/lovego/logc?status.svg)](https://godoc.org/github.com/lovego/logc)

# Install

```
sudo wget -O /usr/local/bin/logc 'https://github.com/lovego/logc/releases/download/170706/logc'
sudo chmod +x /usr/local/bin/logc
```

# Usage
```
logc <your_logc_config_file.yml>
```
See <a href="testdata/logc.yml">logc.yml</a> for full config format.


