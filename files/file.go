package files

import (
	"encoding/json"
	"log"
	"os"

	"github.com/lovego/logc/logd"
	"github.com/lovego/xiaomei/utils/fs"
)

type File struct {
	org, name, path, offsetPath string
	file                        *os.File
	reader                      *json.Decoder
	logger                      *log.Logger
	logd                        *logd.Logd
}

func New(org, name, path string, logd *logd.Logd) *File {
	if fs.NotExist(path) {
		return nil
	}
	f := &File{org: org, name: name, path: path, logd: logd,
		offsetPath: `logc/` + name + `/` + name + `.offset`,
	}
	if !f.openFiles() {
		return nil
	}
	f.reader = json.NewDecoder(f.file)
	return f
}

func (f *File) openFiles() bool {
	dir := `logc/` + f.name
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("mkdir %s error: %v\n", dir, err)
		return false
	}
	if logFile, err := fs.OpenAppend(dir + `/` + f.name + `.log`); err == nil {
		f.logger = log.New(logFile, ``, log.LstdFlags)
	} else {
		log.Printf("open %s: %v\n", f.path, err)
		return false
	}
	var err error
	if f.file, err = os.Open(f.path); err != nil {
		f.logger.Printf("open %s: %v\n", f.path, err)
		return false
	}
	return true
}
