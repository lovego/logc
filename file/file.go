package file

import (
	"encoding/json"
	"os"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type File struct {
	org, name, path string
	logFile, file   *os.File
	reader          *json.Decoder
}

func New(org, name, path string) *File {
	if fs.NotExist(path) {
		return nil
	}
	f := &File{org: org, name: name, path: path}
	f.openLog()
	var err error
	if f.file, err = os.Open(f.path); err != nil {
		f.log(`open ` + f.path + `: ` + err.Error())
		return nil
	}
	f.reader = json.NewDecoder(f.file)
	return f
}

func (f *File) openLog() {
	dir := `logc/` + f.name
	var err error
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	f.logFile = fs.OpenAppend(dir + `/` + f.name + `.log`)
}

func (f *File) log(msg string) {
	utils.Logf(f.logFile, msg)
}
