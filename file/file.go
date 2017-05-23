package file

import (
	"encoding/json"
	"os"

	"github.com/lovego/xiaomei/utils"
)

type File struct {
	org, name, path string
	file, log       *os.File
	reader          *json.Decoder
}

func New(org, name, path string) *File {
	f := &File{org: org, name: name, path: path}
	if f.file = f.openFile(); f.file == nil {
		return nil
	}
	if f.log = f.openLog(); f.log == nil {
		return nil
	}
	f.reader = json.NewDecoder(f.file)
	return f
}

func (f *File) openFile() *os.File {
	file, err := os.Open(f.path)
	if err != nil {
		utils.Log(`open ` + f.path + `: ` + err.Error())
		return nil
	}
	if offset := f.readOffset(); offset > 0 {
		if _, err := file.Seek(offset, os.SEEK_SET); err != nil {
			utils.Log(`seek ` + f.path + `: ` + err.Error())
		}
	}
	return file
}

func (f *File) openLog() *os.File {
	dir := `logc/` + f.name
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	log, err := os.Open(dir + `/` + f.name + `.log`)
	if err != nil {
		panic(err)
	}
	return log
}
