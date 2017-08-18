package main

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/watch"
	"github.com/robfig/cron"
)

func main() {
	conf := config.Get()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	startRotate(conf.RotateTime, conf.RotateCmd)
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
			log.Println(buf.String(), err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	go c.Start()
}
