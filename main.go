package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	filepkg "github.com/lovego/logc/file"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/httputil"
)

const defaultLogdAddr = `192.168.202.12:30432`

func main() {
	org, logdAddr := getParams()
	if logdAddr == `` {
		logdAddr = defaultLogdAddr
	}
	utils.Log(`logc starting. (logd: ` + logdAddr + `)`)
	listenOrgFiles(org, logdAddr)
	// select {}
}

func listenOrgFiles(org, logdAddr string) {
	files := map[string]string{}
	httputil.Http(http.MethodGet, `http://`+logdAddr+`/files?org=`+org, nil, nil, &files)
	wg := sync.WaitGroup{}
	for name, path := range files {
		if file := filepkg.New(org, name, path); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.Listen()
			}()
		}
	}
	wg.Wait()
}

func getParams() (org, logdAddr string) {
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.Usage = usage
	help := flagset.Bool(`help`, false, `print usage info.`)
	flagset.Parse(os.Args[1:])
	args := flagset.Args()

	if len(args) == 0 || len(args) > 2 || *help {
		usage()
		os.Exit(1)
	}
	org = args[0]
	if len(args) > 1 {
		logdAddr = args[1]
	}
	return
}

func usage() {
	fmt.Printf(`a client which listen files, collect contents, and push to logd server
Usage:
  logc <org> [logd-address]
  default address: %s
  example: logc data-visual
`, defaultLogdAddr)
}
