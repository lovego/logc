package config

import (
	"encoding/json"
	"log"
	"path/filepath"
	"time"
)

type Config struct {
	Name              string                 `yaml:"name"`
	Elasticsearch     []string               `yaml:"elasticsearch"`
	MergeJson         string                 `yaml:"-"`
	MergeData         map[string]interface{} `yaml:"mergeData"`
	BatchSize         int                    `yaml:"batchSize"`
	BatchWait         string                 `yaml:"batchWait"`
	BatchWaitDuration time.Duration          `yaml:"-"`
	RotateTime        string                 `yaml:"rotateTime"`
	RotateCmd         []string               `yaml:"rotateCmd"`
	Mailer            string                 `yaml:"mailer"`
	Keepers           []string               `yaml:"keepers"`
	Files             []*File                `yaml:"files"`
}

type File struct {
	Path    string                            `yaml:"path"`
	Index   string                            `yaml:"index"`
	Type    string                            `yaml:"type"`
	Mapping map[string]map[string]interface{} `yaml:"mapping"`
}

func check(conf *Config) {
	checkMergeData(conf)
	checkBatchWait(conf)
	for _, file := range conf.Files {
		checkFile(file)
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

func checkBatchWait(conf *Config) {
	if conf.BatchWait == `` {
		conf.BatchWaitDuration = -1
		return
	}
	duration, err := time.ParseDuration(conf.BatchWait)
	if err != nil {
		log.Fatalf("parse batchWait error: %v", err)
	}
	conf.BatchWaitDuration = duration
}

func checkFile(file *File) {
	if file.Index == `` {
		log.Fatalf("index missing for file: %+v", file)
	}
	if file.Type == `` {
		log.Fatalf("type missing for file: %+v", file)
	}
	if file.Path == `` {
		log.Fatalf("path missing for file: %+v", file)
	} else {
		file.Path = filepath.Clean(file.Path)
	}
}
