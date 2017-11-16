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
		flag.CommandLine.Usage = usage
		flag.Parse()

		if flag.NArg() != 1 {
			usage()
			os.Exit(1)
		}

		configFile := flag.Arg(0)
		if notExist(configFile) {
			log.Fatal(configFile + ` not exist.`)
		}
		conf := parse(configFile)
		conf.check()
		theConf = conf
	}
	return theConf
}

func parse(configFile string) Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	confs := struct {
		Envs map[string]Config `yaml:"envs"`
	}{}
	if err := yaml.Unmarshal(content, &confs); err != nil {
		log.Fatal(err)
	}
	env := os.Getenv(`GOENV`)
	if env == `` {
		env = `dev`
	}
	conf, ok := confs.Envs[env]
	if !ok {
		log.Fatalf(`%s: %s: undefined.`, configFile, env)
	}
	return conf
}

func notExist(p string) bool {
	_, err := os.Stat(p)
	return err != nil && os.IsNotExist(err)
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`logc watch files, collect content, and push to logd server. (version: 17.9.26)
usage: %s <yaml-config-file>
`, os.Args[0])
	flag.PrintDefaults()
}
