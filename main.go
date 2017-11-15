package main

import (
	"bytes"
	logpkg "log"
	"os/exec"
	"time"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/watch"
	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/logger"
	"github.com/robfig/cron"
)

func main() {
	conf := config.Get()
	logpkg.Printf(
		"logc starting. (log-es: %v, merge: %v)\n",
		conf.ElasticSearch, conf.MergeData,
	)

	collector.SetAlarm(theAlarm)
	collector.SetLogger(log)
	reader.SetBatch(conf.BatchSize, conf.BatchWaitDuration)

	startRotate(conf.RotateTime, conf.RotateCmd, log)

	watchFiles(conf, log)
}

func watchFiles(conf config.Config, log *logger.Logger) {
	files := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		files[file.Path] = collectorGetter(
			file.Path, pusher.NewGetter(conf.ElasticSearch, file, conf.MergeJson),
		)
	}
	watch.Watch(files)
}

func collectorGetter(
	path string, pusherGetter collector.PusherGetter,
) func() watch.Collector {
	return func() watch.Collector {
		if c := collector.New(path, pusherGetter); c == nil {
			return nil // must
		} else {
			return c // nil pointer makes a non nil interface
		}
	}
}

func startRotate(timeSpec string, cmd []string, log *logger.Logger) {
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
			log.Errorf("rotate failed: %s, %v", buf.String(), err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	go c.Start()
}
