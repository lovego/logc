package collector

import (
	"os"

	"github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/outputs"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

type Reader interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Close()
}

type Collector struct {
	logFile    *os.File
	logger     *loggerpkg.Logger
	reader     Reader
	output     outputs.Output
	writeEvent chan struct{}
	closeEvent chan struct{}
}

func New(path string, output outputs.Output) *Collector {
	var file, logFile *os.File
	if file = openFile(path); file != nil {
		if logcPath := getLogcPath(path, file); logcPath != `` {
			if logFile := openLogFile(logcPath + `.log`); logFile != nil {
				theLogger := loggerpkg.New(``, logFile, theAlarm)
				theLogger.Printf("collect %s", path)
				if theReader := reader.New(file, logcPath+`.offset`, theLogger); theReader != nil {
					collector := &Collector{
						logFile:    logFile,
						logger:     theLogger,
						reader:     theReader,
						output:     output,
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
