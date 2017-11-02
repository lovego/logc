package main

import (
	"bytes"
	logpkg "log"
	"os"
	"os/exec"
	"time"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/watch"
	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/logger"
	"github.com/lovego/xiaomei/utils/mailer"
	"github.com/robfig/cron"
)

func main() {
	conf := config.Get()
	logpkg.Printf(
		"logc starting. (log-es: %v, merge: %v)\n",
		conf.Elasticsearch, conf.MergeData,
	)

	theAlarm := getAlarm(conf.Name, conf.Mailer, conf.Keepers)
	log := logger.New(``, os.Stderr, theAlarm)

	collector.SetAlarm(theAlarm)
	collector.SetLogger(log)
	reader.SetBatch(conf.BatchSize, conf.BatchWaitDuration)

	pusher.CreateMappings(conf.Elasticsearch, conf.Files, log)
	startRotate(conf.RotateTime, conf.RotateCmd, log)

	watchFiles(conf)
}

func watchFiles(conf config.Config) {
	files := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		files[file.Path] = collectorGetter(
			file.Path, pusher.NewGetter(conf.Elasticsearch, file.Index, file.Type, conf.MergeJson),
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

func getAlarm(name, mailerUrl string, keepers []string) *alarm.Alarm {
	m, err := mailer.New(mailerUrl)
	if err != nil {
		logpkg.Panic(err)
	}
	env := os.Getenv(`GOENV`)
	if env == `` {
		env = `dev`
	}
	return alarm.New(
		name+`_`+env+`_logc`, alarm.MailSender{Receivers: keepers, Mailer: m},
		0, 5*time.Second, 30*time.Second,
	)
}
