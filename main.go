package main

import (
	"log"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/watch"
)

func main() {
	conf := config.Get()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	pusher.CreateMappings(conf.LogdAddr, conf.Files)

	files := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		files[file.Path] = collectorGetter(file.Path,
			pusher.NewGetter(conf.LogdAddr, file.Org, file.Name, conf.MergeJson),
		)
	}
	watch.Watch(files)
}

func collectorGetter(path string, pusherGetter collector.PusherGetter) func() watch.Collector {
	return func() watch.Collector {
		if c := collector.New(path, pusherGetter); c == nil {
			return nil // must
		} else {
			return c // nil pointer makes a non nil interface
		}
	}
}
