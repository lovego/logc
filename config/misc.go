package config

import (
	"log"
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/logger"
	"github.com/lovego/xiaomei/utils/mailer"
)

var theEnv string

func Env() string {
	if theEnv == `` {
		theEnv = os.Getenv(`GOENV`)
		if theEnv == `` {
			theEnv = `dev`
		}
	}
	return theEnv
}

var theAlarm *alarm.Alarm

func Alarm() *alarm.Alarm {
	if theAlarm == nil {
		conf := Get()
		m, err := mailer.New(conf.Mailer)
		if err != nil {
			log.Panic(err)
		}
		env := os.Getenv(`GOENV`)
		if env == `` {
			env = `dev`
		}
		theAlarm = alarm.New(
			conf.Name+`_`+env+`_logc`, alarm.MailSender{Receivers: conf.Keepers, Mailer: m},
			0, 5*time.Second, 30*time.Second,
		)
	}
	return theAlarm
}

var theLogger *logger.Logger

func Logger() *logger.Logger {
	if theLogger == nil {
		theLogger = logger.New(``, os.Stderr, theAlarm)
	}
	return theLogger
}
