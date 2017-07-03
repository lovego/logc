package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	filespkg "github.com/lovego/logc/files"
	logdpkg "github.com/lovego/logc/logd"
	"gopkg.in/yaml.v2"
)

type File struct {
	Name    string                            `yaml:"name"`
	Path    string                            `yaml:"path"`
	Mapping map[string]map[string]interface{} `yaml:"mapping"`
}

func main() {
	logdAddr, mergeJson, orgName, filesAry := getParams()
	log.Printf(
		"logc starting. (logd: %s, merge: %s, org: %s)\n",
		logdAddr, mergeJson, orgName,
	)
	logd, err := logdpkg.New(logdAddr, mergeJson)
	if err != nil {
		log.Fatal(err)
	}
	createOrgFiles(logd, orgName, filesAry)
	listenOrgFiles(logd, orgName, filesAry)
}

func createOrgFiles(logd *logdpkg.Logd, orgName string, filesAry []File) {
	var filesMapping []map[string]interface{}
	for _, file := range filesAry {
		filesMapping = append(filesMapping, map[string]interface{}{
			`name`: file.Name, `mapping`: file.Mapping,
		})
	}
	if err := logd.Create(orgName, filesMapping); err != nil {
		log.Fatal(err)
	}
}

func listenOrgFiles(logd *logdpkg.Logd, orgName string, filesAry []File) {
	wg := sync.WaitGroup{}
	for _, info := range filesAry {
		if file := filespkg.New(orgName, info.Name, info.Path, logd); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.Listen()
			}()
		}
	}
	wg.Wait()
}

func getParams() (logd, merge, org string, files []File) {
	flag.StringVar(&merge, `merge`, ``, "merge the `json` object into data lines")
	help := flag.Bool(`help`, false, `print help message.`)
	flag.CommandLine.Usage = usage
	flag.Parse()

	if flag.NArg() != 3 || *help {
		usage()
		os.Exit(1)
	}
	logd = flag.Arg(0)
	org = flag.Arg(1)
	files = parseFiles(flag.Arg(2))
	return
}

func parseFiles(confFile string) []File {
	content, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := struct {
		Files []File
	}{}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Fatal(err)
	}
	return conf.Files
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s listen files, and push content to logd server.\n\n"+
		"Usage: %s [options] logd-addr org-name conf-file\n"+
		"Options:\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
}
