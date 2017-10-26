package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Get() Config {
	configFile := flag.Arg(0)
	if notExist(configFile) {
		log.Fatal(configFile + ` not exist.`)
	}
	conf := parse(configFile)
	check(&conf)
	return conf
}

func init() {
	flag.CommandLine.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`logc watch files, collect content, and push to logd server. (version: 17.9.26)
usage: %s <yaml-config-file>
`, os.Args[0])
	flag.PrintDefaults()
}

func parse(configFile string) Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	envConfs := map[string]Config{}
	if err := yaml.Unmarshal(content, &envConfs); err != nil {
		log.Fatal(err)
	}
	for _, envConf := range envConfs {
		for _, file := range envConf.Files {
			cleanFileMapping(file.Mapping)
		}
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

func notExist(p string) bool {
	_, err := os.Stat(p)
	return err != nil && os.IsNotExist(err)
}

func cleanFileMapping(mapping map[string]map[string]interface{}) {
	for _, value := range mapping {
		for k, v := range value {
			value[k] = cleanupMapValue(v)
		}
	}
}

func cleanupMapValue(value interface{}) interface{} {
	if data, ok := value.(map[interface{}]interface{}); ok {
		res := make(map[string]interface{})
		for k, v := range data {
			res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
		}
		return res
	}
	return value
}
