package collector

import (
	"os"

	readerpkg "github.com/lovego/logc/collector/reader"
	"github.com/lovego/logc/outputs"
	loggerpkg "github.com/lovego/logger"
)

type Reader interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Close()
}

type Collector struct {
	id         string
	logFile    *os.File
	logger     *loggerpkg.Logger
	reader     Reader
	output     outputs.Output
	writeEvent chan struct{}
	closeEvent chan struct{}
}

func New(path, name string, outputMaker func(*loggerpkg.Logger) outputs.Output) *Collector {
	var file, logFile *os.File
	if file = openFile(path); file != nil {
		if logcPath := getLogcPath(path, name, file); logcPath != `` {
			if logFile := openLogFile(logcPath + `.log`); logFile != nil {
				logger := loggerpkg.New(logFile)
				logger.SetAlarm(theAlarm)
				logger.Infof("collect %s", path)

				collectorId := path + `:` + name
				if reader := readerpkg.New(collectorId, file, logcPath+`.offset`, logger); reader != nil {
					if output := outputMaker(logger); output != nil {
						collector := &Collector{
							id:         collectorId,
							logFile:    logFile,
							logger:     logger,
							reader:     reader,
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
	}
	freeResource(file, logFile)
	return nil
}
