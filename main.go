package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	filespkg "github.com/lovego/logc/files"
	logdpkg "github.com/lovego/logc/logd"
	"github.com/lovego/xiaomei/utils"
)

func main() {
	logdAddr, mergeJson, orgName := getParams()
	utils.Logf(
		`logc starting. (logd address: %s, merge json: %s, org name: %s)`,
		logdAddr, mergeJson, orgName,
	)
	logd := logdpkg.New(logdAddr, mergeJson)
	if logd == nil {
		os.Exit(1)
	}
	listenOrgFiles(logd, orgName)
}

func listenOrgFiles(logd *logdpkg.Logd, orgName string) {
	filesAry := logd.FilesOf(orgName)
	wg := sync.WaitGroup{}
	for _, info := range filesAry {
		if file := filespkg.New(orgName, info[`name`], info[`path`], logd); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.Listen()
			}()
		}
	}
	wg.Wait()
}

func getParams() (logd, merge, org string) {
	flag.StringVar(&merge, `merge`, ``, "merge the `json` object into data lines")
	help := flag.Bool(`help`, false, `print help message.`)
	flag.CommandLine.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 || *help {
		usage()
		os.Exit(1)
	}
	logd = flag.Arg(0)
	org = flag.Arg(1)
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s listen files, and push content to logd server.\n\n"+
		"Usage: %s [options] logd-addr org-name\n"+
		"Options:\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
}
