package config

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
)

type Config struct {
	LogdAddr   string `yaml:"logdAddr"`
	MergeJson  string
	MergeData  map[string]interface{} `yaml:"mergeData"`
	Files      []*File                `yaml:"files"`
	BatchSize  int                    `yaml:"batchSize"`
	RotateTime string                 `yaml:"rotateTime"`
	RotateCmd  []string               `yaml:"rotateCmd"`
	Name       string                 `yaml:"name"`
	Mailer     string                 `yaml:"mailer"`
	Keepers    []string               `yaml:"keepers"`
}

type File struct {
	Org     string                            `yaml:"org"`
	Name    string                            `yaml:"name"`
	Path    string                            `yaml:"path"`
	Mapping map[string]map[string]interface{} `yaml:"mapping"`
}

func check(conf *Config) {
	checkLogdAddress(conf)
	checkMergeData(conf)
	if len(conf.Files) == 0 {
		log.Fatal(`files required.`)
	}
	for _, file := range conf.Files {
		checkFile(file)
	}
}

func checkLogdAddress(conf *Config) {
	addr := conf.LogdAddr
	if addr == `` {
		log.Fatal(`logd address required.`)
	}
	if !strings.HasPrefix(addr, `http://`) && !strings.HasPrefix(addr, `https://`) {
		conf.LogdAddr = `http://` + addr
	}
}

func checkMergeData(conf *Config) {
	if len(conf.MergeData) > 0 {
		if buf, err := json.Marshal(conf.MergeData); err != nil {
			log.Fatalf("marshal merge data: %v", err)
		} else {
			conf.MergeJson = string(buf)
		}
	}
}

func checkFile(file *File) {
	if file.Org == `` {
		log.Fatalf("org missing for file: %+v", file)
	}
	if file.Name == `` {
		log.Fatalf("name missing for file: %+v", file)
	}
	if file.Path == `` {
		log.Fatalf("path missing for file: %+v", file)
	} else {
		file.Path = filepath.Clean(file.Path)
	}
}
