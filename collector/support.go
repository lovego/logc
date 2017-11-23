package collector

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/lovego/xiaomei/utils/alarm"
	"github.com/lovego/xiaomei/utils/fs"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

var logger *loggerpkg.Logger
var theAlarm *alarm.Alarm

func Setup(l *loggerpkg.Logger, a *alarm.Alarm) {
	logger = l
	theAlarm = a
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Errorln(err) // os.PathError is enough
		}
		return nil
	}
	return file
}

func getLogcPath(path, collectorId string, f *os.File) string {
	fi, err := f.Stat()
	if err != nil {
		logger.Errorln("stat:", err)
		return ``
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || sys == nil {
		logger.Errorf("unexpected FileInfo.Sys(): %#v", fi.Sys())
		return ``
	}
	ino := strconv.FormatUint(sys.Ino, 10)
	return filepath.Join(filepath.Dir(path), `.logc`, collectorId+`.`+ino)
}

func openLogFile(path string) *os.File {
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		logger.Errorf("logger: %v", err) // os.PathError is enough
		return nil
	}
	if file, err := fs.OpenAppend(path); err == nil {
		return file
	} else {
		logger.Errorf("logger: %v", err) // os.PathError is enough
		return nil
	}
}

func freeResource(file, logFile *os.File) {
	if file != nil {
		if err := file.Close(); err != nil {
			logger.Errorf("close error: %v", err)
		}
	}
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			logger.Errorf("logger: close error: %v", err)
		}
	}
}
