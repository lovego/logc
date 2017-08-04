package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/source"
	"github.com/lovego/logc/watch"
	"github.com/lovego/xiaomei/utils/fs"
)

func main() {
	conf := config.Get()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	pusher.CreateMappings(conf.LogdAddr, conf.Files)

	collectors := make(map[string]watch.Collector)
	for _, file := range conf.Files {
		collectors[file.Path] = getCollector(file, conf.LogdAddr, conf.MergeJson)
	}
	watch.Watch(collectors)
}

func getCollector(file *config.File, logdAddr, mergeJson string) watch.Collector {
	keyPath := filepath.Join(`logc`, file.Org, file.Name)
	logger := getLogger(keyPath + `.log`)

	return collector.New(
		file.Path,
		source.New(file.Path, keyPath+`.offset`, logger),
		pusher.New(logdAddr, file.Org, file.Name, mergeJson, logger),
		logger,
	)
}

func getLogger(path string) *log.Logger {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("mkdir %s error: %v\n", dir, err)
	}
	if logFile, err := fs.OpenAppend(path); err == nil {
		return log.New(logFile, ``, log.LstdFlags)
	} else {
		log.Fatal("open %s: %v\n", path, err)
		return nil
	}
}
