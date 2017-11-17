package config

import (
	"log"
	"path/filepath"

	"github.com/lovego/logc/collector/reader"
)

type Config struct {
	Name    string       `yaml:"name"`
	Mailer  string       `yaml:"mailer"`
	Keepers []string     `yaml:"keepers"`
	Batch   reader.Batch `yaml:"batch"`
	Rotate  Rotate       `yaml:"rotate"`
	Files   []File       `yaml:"files"`
}

type Rotate struct {
	Time string   `yaml:"time"`
	Cmd  []string `yaml:"cmd"`
}

type File struct {
	Path    string                   `yaml:"path"`
	Outputs []map[string]interface{} `yaml:"outputs"`
}

func (conf *Config) check() {
	if conf.Name == `` {
		log.Fatal("config: empty name")
	}
	for i, file := range conf.Files {
		if file.Path == `` {
			log.Fatalf("path missing for file: %+v", file)
		} else {
			file.Path = filepath.Clean(file.Path)
		}
		conf.Files[i] = file
	}
}
