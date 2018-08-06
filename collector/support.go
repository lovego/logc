package collector

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/lovego/alarm"
	"github.com/lovego/fs"
	loggerpkg "github.com/lovego/logger"
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
			logger.Error(err) // os.PathError is enough
		}
		return nil
	}
	return file
}

func getLogcPath(path, name string, f *os.File) string {
	fi, err := f.Stat()
	if err != nil {
		logger.Error("stat:", err)
		return ``
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || sys == nil {
		logger.Errorf("unexpected FileInfo.Sys(): %#v", fi.Sys())
		return ``
	}
	ino := strconv.FormatUint(sys.Ino, 10)
	return filepath.Join(filepath.Dir(path), `.logc`, name+`.`+ino)
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
			logger.Errorf("close logger error: %v", err)
		}
	}
}
