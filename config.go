package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogdAddr  string                 `yaml:"logdAddr"`
	MergeData map[string]interface{} `yaml:"mergeData"`
	Files     []File                 `yaml:"files"`
}

type File struct {
	OrgName string                            `yaml:"orgName"`
	Name    string                            `yaml:"name"`
	Path    string                            `yaml:"path"`
	Mapping map[string]map[string]interface{} `yaml:"mapping"`
}

func getConfig() Config {
	help := flag.Bool(`help`, false, `print help message.`)
	flag.CommandLine.Usage = usage
	flag.Parse()

	if flag.NArg() > 1 || *help {
		usage()
		os.Exit(1)
	}
	configFile := flag.Arg(0)
	if configFile == `` {
		configFile = `logc.yml`
	}
	if notExist(configFile) {
		log.Fatal(configFile + ` not exist.`)
	}
	return parseConfig(configFile)
}

func parseConfig(configFile string) Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	envConfs := map[string]Config{}
	if err := yaml.Unmarshal(content, &envConfs); err != nil {
		log.Fatal(err)
	}
	env := os.Getenv(`GOENV`)
	if env == `` {
		env = `dev`
	}
	conf, ok := envConfs[env]
	if !ok {
		log.Fatalf(`%s: %s: undefined.`, configFile, env)
	}
	return conf
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [config-file] (default: logc.yml)\n"+
			"logc listen files, and push content to logd server.\n", os.Args[0],
	)
}

func notExist(p string) bool {
	_, err := os.Stat(p)
	return err != nil && os.IsNotExist(err)
}
