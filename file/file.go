package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type File struct {
	org, name, path, offsetPath string
	log, file                   *os.File
	reader                      *json.Decoder
}

func New(org, name, path string) *File {
	if fs.NotExist(path) {
		return nil
	}
	f := &File{org: org, name: name, path: path,
		offsetPath: `logc/` + name + `/` + name + `.offset`,
	}
	f.openLog()
	var err error
	if f.file, err = os.Open(f.path); err != nil {
		f.Log(`open %s: %v`, f.path, err)
		return nil
	}
	f.Log(`collect ` + path)
	f.seekToSavedOffset()
	f.reader = json.NewDecoder(f.file)
	return f
}

func (f *File) openLog() {
	dir := `logc/` + f.name
	var err error
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	f.log = fs.OpenAppend(dir + `/` + f.name + `.log`)
}

func (f *File) Log(format string, args ...interface{}) {
	utils.Logf(f.log, fmt.Sprintf(format, args...))
}
