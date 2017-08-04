package main

import (
	"log"
	"sync"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/source"
	"github.com/lovego/logc/watch"
)

func main() {
	conf := getConfig()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	pusher.CreateMappings(conf.LogdAddr, conf.Files)

	collectors := make(map[string]watch.Collector)
	for _, file := range conf.Files {
		collectors[file.Path] = makeCollector(conf.logdAddr, conf.MergeJson, file)
	}
	watch.Watch(collectors)
}

func getCollector(logdAddr, mergeJson, file *config.File) watch.Collector {
	keyPath := filepath.Join(`logc`, file.Org, file.Name)
	logger := getLogger(keyPath + `.log`)

	return collector.New(
		file.Path,
		source.New(file.Path, keyPath+`.offset`, logger),
		pusher.New(logdAddr, file.Org, file.Name, conf.MergeJson, logger),
		logger,
	)
}

func getLogger(dir, name string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("mkdir %s error: %v\n", dir, err)
	}
	logPath := filepath.Join(dir, name+`.log`)
	if logFile, err := fs.OpenAppend(logPath); err == nil {
		c.logger = log.New(logFile, ``, log.LstdFlags)
	} else {
		log.Fatal("open %s: %v\n", logPath, err)
	}
}
