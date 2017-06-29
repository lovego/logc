package files

import (
	"encoding/json"
	"os"

	"github.com/lovego/logc/logd"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type File struct {
	org, name, path, offsetPath string
	log, file                   *os.File
	reader                      *json.Decoder
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
		utils.Logf(`mkdir %s error: %v`, dir, err)
		return false
	}
	var err error
	if f.log, err = fs.OpenAppend(dir + `/` + f.name + `.log`); err != nil {
		f.Log(`open %s: %v`, f.path, err)
		return false
	}
	if f.file, err = os.Open(f.path); err != nil {
		f.Log(`open %s: %v`, f.path, err)
		return false
	}
	return true
}

func (f *File) Log(format string, args ...interface{}) {
	utils.FLogf(f.log, format, args...)
}
