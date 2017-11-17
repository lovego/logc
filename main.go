package main

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/watch"
	"github.com/robfig/cron"
)

var logger = config.Logger()

func main() {
	conf := config.Get()
	reader.Setup(conf.Batch, logger)
	collector.Setup(logger, config.Alarm())

	files := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		files[file.Path] = collectorGetter(file.Path, pusher.NewGetter(file))
	}

	startRotate(conf.Rotate.Time, conf.Rotate.Cmd)

	log.Printf("logc starting. (log-es: %v)\n", conf.ElasticSearch)
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

func startRotate(timeSpec string, cmd []string) {
	if timeSpec == `` || len(cmd) == 0 {
		return
	}
	c := cron.New()
	err := c.AddFunc(timeSpec, func() {
		var buf bytes.Buffer
		cmd := exec.Command(cmd[0], cmd[1:]...)
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		if err := cmd.Run(); err != nil {
			logger.Errorf("rotate failed: %s, %v", buf.String(), err)
		}
	})
	if err != nil {
		logger.Error(err)
	}
	go c.Start()
}
