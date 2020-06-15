package config

import (
	"log"
	"os"
	"time"

	"github.com/lovego/alarm"
	"github.com/lovego/logger"
	"github.com/lovego/mailer"
)

var theAlarm *alarm.Alarm

func Alarm() *alarm.Alarm {
	if theAlarm == nil {
		conf := Get()
		m, err := mailer.New(conf.Mailer)
		if err != nil {
			log.Panic(err)
		}
		theAlarm = alarm.New(
			alarm.MailSender{Receivers: conf.Keepers, Mailer: m},
			0, 5*time.Second, 30*time.Second, alarm.SetPrefix(conf.Name+`_logc`),
		)
	}
	return theAlarm
}

var theLogger *logger.Logger

func Logger() *logger.Logger {
	if theLogger == nil {
		theLogger = logger.New(os.Stderr)
		theLogger.SetAlarm(Alarm())
	}
	return theLogger
}
