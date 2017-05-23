package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	filepkg "github.com/lovego/logc/file"
	"github.com/lovego/xiaomei/utils/httputil"
)

const defaultAddr = `192.168.202.12:30432`

var remoteAddr string

func main() {
	orgName, addr := getParams()
	remoteAddr = addr
	if addr == `` {
		remoteAddr = defaultAddr
	}
	fmt.Println(`remote address: `, remoteAddr)
	listenOrgFiles(orgName)
	// select {}
}

func listenOrgFiles(orgName string) {
	paths := []string{}
	httputil.Http(http.MethodGet, `http://`+path.Join(remoteAddr, `files?org=`+orgName), nil, nil, &paths)
	wg := sync.WaitGroup{}
	for _, p := range paths {
		if file := filepkg.New(orgName, p); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.listen()
			}()
		}
	}
	wg.Wait()
}

func getParams() (orgName, remoteAddr string) {
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.Usage = usage
	help := flagset.Bool(`help`, false, `print usage info.`)
	flagset.Parse(os.Args[1:])
	args := flagset.Args()

	if len(args) == 0 || len(args) > 2 || *help {
		usage()
		os.Exit(1)
	}
	orgName = args[0]
	if len(args) > 1 {
		remoteAddr = args[1]
	}
	return
}

func usage() {
	fmt.Printf(`a client which listen files, collect contents, and push to server
Usage:
  logc <org> [address]
  default address: %s
  example: logc data-visual
`, defaultAddr)
}
