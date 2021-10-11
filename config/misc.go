package config

import (
	"log"
	"os"

	"github.com/lovego/alarm"
	"github.com/lovego/email"
	"github.com/lovego/logger"
)

var theAlarm *alarm.Alarm

func Alarm() *alarm.Alarm {
	if theAlarm == nil {
		conf := Get()
		m, err := email.NewClient(conf.Mailer)
		if err != nil {
			log.Panic(err)
		}
		theAlarm = alarm.New(
			alarm.MailSender{Receivers: conf.Keepers, Mailer: m},
			nil,
			alarm.SetPrefix(conf.Name+`_logc`),
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
