package collector

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/fs"
)

var log = config.Logger()
var theAlarm = config.Alarm()

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorln(err) // os.PathError is enough
		}
		return nil
	}
	return file
}

func getLogcPath(path string, f *os.File) string {
	fi, err := f.Stat()
	if err != nil {
		log.Errorln("stat:", err)
		return ``
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || sys == nil {
		log.Errorf("unexpected FileInfo.Sys(): %#v", fi.Sys())
		return ``
	}
	return filepath.Join(filepath.Dir(path), `logc`, strconv.FormatUint(sys.Ino, 10))
}

func openLogFile(path string) *os.File {
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		log.Errorf("logger: %v", err) // os.PathError is enough
		return nil
	}
	if file, err := fs.OpenAppend(path); err == nil {
		return file
	} else {
		log.Errorf("logger: %v", err) // os.PathError is enough
		return nil
	}
}

func freeResource(file, logFile *os.File) {
	if file != nil {
		if err := file.Close(); err != nil {
			log.Errorf("close error: %v", err)
		}
	}
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			log.Errorf("logger: close error: %v", err)
		}
	}
}
