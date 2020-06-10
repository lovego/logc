package main

import (
	"bytes"
	"os/exec"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/outputs"
	"github.com/lovego/logc/watch"
	"github.com/robfig/cron"
)

var logger = config.Logger()

func main() {
	conf := config.Get()
	reader.Setup(conf.Config, logger)
	collector.Setup(logger, config.Alarm())
	outputs.Setup(logger)

	startRotate(conf.Rotate.Time, conf.Rotate.Cmd)

	collectorMakers := getCollectorMakers(conf.Files)

	watch.Watch(collectorMakers)
}

func getCollectorMakers(files map[string]map[string]map[string]interface{}) map[string]func() []watch.Collector {
	makers := make(map[string]func() []watch.Collector)
	for path, collectorsConf := range files {
		makers[path] = getCollectorsMaker(path, collectorsConf)
	}
	return makers
}

func getCollectorsMaker(path string, collectorsConf map[string]map[string]interface{}) func() []watch.Collector {
	return func() (collectors []watch.Collector) {
		for name, outputConf := range collectorsConf {
			if outputMaker := outputs.Maker(path+`:`+name, outputConf); outputMaker != nil {
				if c := collector.New(path, name, outputMaker); c != nil {
					collectors = append(collectors, c)
				}
			}
		}
		return collectors
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
