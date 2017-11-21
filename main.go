package main

import (
	"bytes"
	"log"
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
	reader.Setup(conf.Batch, logger)
	collector.Setup(logger, config.Alarm())
	outputs.Setup(logger)

	startRotate(conf.Rotate.Time, conf.Rotate.Cmd)

	collectorMakers := getCollectorMakers(conf.Files)

	log.Printf("logc starting.\n")
	watch.Watch(collectorMakers)
}

func getCollectorMakers(files []config.File) map[string]func() []watch.Collector {
	makers := make(map[string]func() []watch.Collector)
	for _, file := range files {
		makers[file.Path] = getCollectorsMaker(file)
	}
	return makers
}

func getCollectorsMaker(file config.File) func() []watch.Collector {
	// outputs.CheckConfig(file.Path, file.Outputs)
	return func() (collectors []watch.Collector) {
		for _, outputConf := range file.Outputs {
			if outputMaker := outputs.Maker(outputConf, file.Path); outputMaker != nil {
				if c := collector.New(file.Path, outputMaker); c != nil {
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
