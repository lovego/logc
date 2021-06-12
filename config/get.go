package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var theConf Config

func Get() Config {
	if theConf.Name == `` {
		configFile := getArguments()
		if notExist(configFile) {
			log.Fatal(configFile + ` not exist.`)
		}
		conf := parse(configFile)
		conf.check()
		conf.markEnv()
		theConf = *conf
	}
	return theConf
}

func parse(configFile string) *Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	conf := Config{}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Fatal(err)
	}
	return &conf
}

func notExist(p string) bool {
	_, err := os.Stat(p)
	return err != nil && os.IsNotExist(err)
}

func getArguments() string {
	flag.CommandLine.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	return flag.Arg(0)
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`logc watch files, collect content, and push to log server, such as elasticsearch etc.

version: 21.6.12
usage: %s <yaml-config-file>
`, os.Args[0])
	flag.PrintDefaults()
}
