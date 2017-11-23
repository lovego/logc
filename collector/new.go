package collector

import (
	"os"

	readerpkg "github.com/lovego/logc/collector/reader"
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

func New(path, collectorId string, outputMaker func(*loggerpkg.Logger) outputs.Output) *Collector {
	var file, logFile *os.File
	if file = openFile(path); file != nil {
		if logcPath := getLogcPath(path, collectorId, file); logcPath != `` {
			if logFile := openLogFile(logcPath + `.log`); logFile != nil {
				logger := loggerpkg.New(``, logFile, theAlarm)
				logger.Printf("collect %s", path)
				if reader := readerpkg.New(file, logcPath+`.offset`, logger); reader != nil {
					if output := outputMaker(logger); output != nil {
						collector := &Collector{
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
