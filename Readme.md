# Logc
A log collector(like logstash, fluentd) written by golang.
Now only support read from files and output to elasticSearch.

[![Build Status](https://github.com/lovego/logc/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/logc/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/logc/badge.svg?branch=master&1)](https://coveralls.io/github/lovego/logc)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/logc)](https://goreportcard.com/report/github.com/lovego/logc)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/logc)](https://pkg.go.dev/github.com/lovego/logc@v0.0.2)

# Install
```
go get github.com/lovego/logc
```

# Usage
```
logc <your_logc_config_file.yml>
```
See <a href="testdata/logc.yml">logc.yml</a> for full config format.


