package file

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
				f.Log(`seekToSavedOffset %s: offset(%d) size(%d)`, f.path, offset, size)
			} else {
				f.Log(`seekToSavedOffset %s error: %v`, f.path, err)
			}
		} else {
			f.Log(`seekToSavedOffset %s: offset(%d) exceeds size(%d)`, f.path, offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (f *File) seekFrontIfTruncated() {
	if offset := f.offset(); offset > 0 {
		if size := f.size(); size > 0 && offset > size {
			if _, err := f.file.Seek(0, os.SEEK_SET); err == nil {
				f.Log(`seekFront (size: %d, offset: %d)`, size, offset)
			} else {
				f.Log(`seekFront (size: %d, offset: %d) error: %v`, size, offset, err)
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
		f.Log(`write offset error: %v`, err.Error())
	}
	return offsetStr
}

func (f *File) readOffset() int64 {
	if fs.NotExist(f.offsetPath) {
		return 0
	}
	content, err := ioutil.ReadFile(f.offsetPath)
	if err != nil {
		f.Log(`read offset %s: %v`, f.offsetPath, err)
		return 0
	}
	offset, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		f.Log(`parse offset %s(%s): %s`, f.offsetPath, content, err)
		return 0
	}
	return offset
}

func (f *File) offset() int64 {
	if offset, err := f.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		f.Log(`get offset error: %v`, err)
		return -1
	}
}

func (f *File) size() int64 {
	if fi, err := f.file.Stat(); err == nil {
		return fi.Size()
	} else {
		f.Log(`get size error: %v`, err)
		return -1
	}
}
