package config

import (
	"log"
	"time"
)

type Config struct {
	Name              string                 `yaml:"name"`
	ElasticSearch     []string               `yaml:"elasticsearch"`
	MergeData         map[string]interface{} `yaml:"mergeData"`
	BatchSize         int                    `yaml:"batchSize"` // <= 0 means unlimted
	BatchWait         string                 `yaml:"batchWait"`
	BatchWaitDuration time.Duration          `yaml:"-"` // <= 0 means don't wait
	RotateTime        string                 `yaml:"rotateTime"`
	RotateCmd         []string               `yaml:"rotateCmd"`
	Mailer            string                 `yaml:"mailer"`
	Keepers           []string               `yaml:"keepers"`
	Files             []File                 `yaml:"files"`
}

func (conf *Config) check() {
	if conf.Name == `` {
		log.Fatal("config: empty name")
	}
	conf.checkBatchWait()
	for i, file := range conf.Files {
		file.check()
		conf.Files[i] = file
	}
}

func (conf *Config) checkBatchWait() {
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
