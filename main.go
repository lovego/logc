package main

import (
	"log"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/logger"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/reader"
	"github.com/lovego/logc/watch"
)

func main() {
	conf := config.Get()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	pusher.CreateMappings(conf.LogdAddr, conf.Files)

	collectors := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		collectors[file.Path] = getCollectorMaker(file, conf.LogdAddr, conf.MergeJson)
	}
	watch.Watch(collectors)
}

func getCollectorMaker(file *config.File, logdAddr, mergeJson string) func() watch.Collector {
	return func() watch.Collector {
		theLogger := logger.New(file.Path)
		if theLogger == nil {
			return nil
		}
		theReader := reader.New(file.Path, theLogger.Get())
		if theReader == nil {
			return nil
		}
		return collector.New(
			file.Path, theReader, theLogger,
			pusher.New(logdAddr, file.Org, file.Name, mergeJson, theLogger.Get()),
		)
	}
}
