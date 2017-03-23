package main

import (
	"flag"
	"fmt"
	"os"
)

const defaultAddr = `192.168.202.12:30432`

var remoteAddr string

func main() {
	flag.Parse()
	params := flag.Args()
	if len(params) < 1 {
		usage()
		os.Exit(1)
	}
	if len(params) > 1 {
		remoteAddr = params[1]
	}
	if remoteAddr == `` {
		remoteAddr = defaultAddr
	}
	FileInfo(params[0])
}

func usage() {
	fmt.Printf(`
a client which listen files, collect contents, and push to server
Usage:
  logc <org> [address]
  default address: %s
  example: logc data-visual
`, defaultAddr)
}
