package files

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/lovego/xiaomei/utils/fs"
)

// 把文件seek到上次保存的位置
func (f *File) seekToSavedOffset() {
	if offset := f.readOffset(); offset > 0 {
		if size := f.size(); size > 0 && offset <= size {
			if _, err := f.file.Seek(offset, os.SEEK_SET); err == nil {
				f.logger.Printf("seekToSavedOffset %s: offset(%d) size(%d)\n", f.path, offset, size)
			} else {
				f.logger.Printf("seekToSavedOffset %s error: %v\n", f.path, err)
			}
		} else {
			f.logger.Printf("seekToSavedOffset %s: offset(%d) exceeds size(%d)\n", f.path, offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (f *File) seekFrontIfTruncated() {
	if offset := f.offset(); offset > 0 {
		if size := f.size(); size > 0 && offset > size {
			if _, err := f.file.Seek(0, os.SEEK_SET); err == nil {
				f.logger.Printf("seekFront (size: %d, offset: %d)\n", size, offset)
			} else {
				f.logger.Printf("seekFront (size: %d, offset: %d) error: %v\n", size, offset, err)
			}
		}
	}
}

func (f *File) writeOffset() string {
	var offsetStr string
	if offset := f.offset(); offset > 0 {
		offsetStr = strconv.FormatInt(offset, 10)
	} else {
		return ``
	}

	if err := ioutil.WriteFile(f.offsetPath, []byte(offsetStr), 0666); err != nil {
		f.logger.Printf("write offset error: %v\n", err.Error())
	}
	return offsetStr
}

func (f *File) readOffset() int64 {
	if fs.NotExist(f.offsetPath) {
		return 0
	}
	content, err := ioutil.ReadFile(f.offsetPath)
	if err != nil {
		f.logger.Printf("read offset %s: %v\n", f.offsetPath, err)
		return 0
	}
	offset, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		f.logger.Printf("parse offset %s(%s): %s\n", f.offsetPath, content, err)
		return 0
	}
	return offset
}

func (f *File) offset() int64 {
	if offset, err := f.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		f.logger.Printf("get offset error: %v\n", err)
		return -1
	}
}

func (f *File) size() int64 {
	if fi, err := f.file.Stat(); err == nil {
		return fi.Size()
	} else {
		f.logger.Printf("get size error: %v\n", err)
		return -1
	}
}
