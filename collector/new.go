package collector

import (
	"os"

	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/xiaomei/utils/logger"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Close()
}

type Reader interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Close()
}

type PusherGetter interface {
	Get(*logger.Logger) Pusher
}

type Pusher interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	logFile    *os.File
	logger     *logger.Logger
	reader     Reader
	pusher     Pusher
	writeEvent chan struct{}
	closeEvent chan struct{}
}

func New(path string, pusherGetter PusherGetter) *Collector {
	var file, logFile *os.File
	if file = openFile(path); file != nil {
		if logcPath := getLogcPath(path, file); logcPath != `` {
			if logFile := openLogFile(logcPath + `.log`); logFile != nil {
				theLogger := logger.New(``, logFile, theAlarm)
				theLogger.Printf("collect %s", path)
				if theReader := reader.New(file, logcPath+`.offset`, theLogger); theReader != nil {
					collector := &Collector{
						logFile:    logFile,
						logger:     theLogger,
						reader:     theReader,
						pusher:     pusherGetter.Get(theLogger),
						writeEvent: make(chan struct{}, 1),
						closeEvent: make(chan struct{}, 1),
					}
					go collector.loop()
					return collector
				}
			}
		}
	}
	freeResource(file, logFile)
	return nil
}
