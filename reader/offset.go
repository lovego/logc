package reader

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lovego/xiaomei/utils/fs"
)

type offsetFile struct {
	path   string
	file   *os.File
	logger *log.Logger
}

func newOffsetFile(path string, logger *log.Logger) *offsetFile {
	path += `.offset`
	o := &offsetFile{path: path, logger: logger}
	if file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		o.file = file
		return o
	} else {
		logger.Printf("offset: open %s error: %v", path, err)
		return nil
	}
}

func (o *offsetFile) read() int64 {
	if fs.NotExist(o.path) {
		return 0
	}
	content, err := ioutil.ReadFile(o.path)
	if err != nil {
		o.logger.Printf("offset: read %s error: %v", o.path, err)
		return 0
	}
	contentStr := strings.TrimSpace(string(content))
	if len(contentStr) == 0 {
		return 0
	}
	offset, err := strconv.ParseInt(contentStr, 10, 64)
	if err != nil {
		o.logger.Printf("offset: parse %s(%s) error: %s", o.path, contentStr, err)
		return 0
	}
	return offset
}

func (o *offsetFile) save(offset int64) string {
	offsetStr := strconv.FormatInt(offset, 10)
	if err := o.file.Truncate(0); err != nil {
		o.logger.Printf("offset: truncate %s error: %v", o.path, err)
	}
	if _, err := o.file.WriteAt([]byte(offsetStr), 0); err != nil {
		o.logger.Printf("offset: write %s error: %v", o.path, err)
	}
	return offsetStr
}

func (o *offsetFile) rename(newPath string) {
	newPath += `.offset`
	if err := os.Rename(o.path, newPath); err == nil {
		o.path = newPath
	} else {
		o.logger.Printf("offset: rename %s to %s error: %v", o.path, newPath, err)
	}
}

func (o *offsetFile) remove() {
	if err := o.file.Close(); err != nil {
		o.logger.Printf("offset: close %s error: %v", o.path, err)
	}
	if err := os.Remove(o.path); err != nil && !os.IsNotExist(err) {
		o.logger.Printf("offset: remove %s error: %v", o.path, err)
	}
}
