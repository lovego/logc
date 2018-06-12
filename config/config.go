package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lovego/logc/collector/reader"
)

type Config struct {
	Name    string                                       `yaml:"name"`
	Mailer  string                                       `yaml:"mailer"`
	Keepers []string                                     `yaml:"keepers"`
	Batch   reader.Batch                                 `yaml:"batch"`
	Rotate  Rotate                                       `yaml:"rotate"`
	Files   map[string]map[string]map[string]interface{} `yaml:"files"`
	// Files:  filePath   output     key     value
}

type Rotate struct {
	Time string   `yaml:"time"`
	Cmd  []string `yaml:"cmd"`
}

func (conf *Config) check() {
	if conf.Name == `` {
		log.Fatal("config: empty name")
	}
	conf.checkFiles()
}

func (conf *Config) checkFiles() {
	for path, collectors := range conf.Files {
		if path == `` {
			log.Fatalf("empty file path for: %+v", path)
		} else {
			if cleanPath := filepath.Clean(path); cleanPath != path {
				delete(conf.Files, path)
				conf.Files[cleanPath] = collectors
			}
		}
		for collectorId, conf := range collectors {
			if len(conf) == 0 {
				log.Fatalf("%s.%s: empty config.", path, collectorId)
			}
		}
	}
}

func (conf *Config) setByEnv() {
	env := os.Getenv(`GOENV`)
	if env == `production` {
		return
	}
	if env == `` {
		env = `dev`
	}
	conf.Name += `_` + env
	for _, file := range conf.Files {
		if file[`es`] != nil {
			file[`es`][`addrs`] = addEnv2EsAddrs(file[`es`][`addrs`], env)
		}
	}
}

func addEnv2EsAddrs(esAddrs interface{}, env string) interface{} {
	if addrs, ok := esAddrs.([]interface{}); ok && len(addrs) > 0 {
		for i, addr := range addrs {
			if address, ok := addr.(string); ok {
				addrs[i] = address + env + `-`
			}
		}
		return addrs
	}
	return esAddrs
}
