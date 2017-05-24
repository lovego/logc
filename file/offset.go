package file

import (
	"fmt"
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
			if _, err := f.file.Seek(offset, os.SEEK_SET); err != nil {
				f.log(`seek ` + f.path + `: ` + err.Error())
			}
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (f *File) seekFrontIfTruncated() {
	if offset := f.offset(); offset > 0 {
		if size := f.size(); size > 0 && offset > size {
			f.file.Seek(0, os.SEEK_SET)
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

	path := f.name + `/` + f.name + `.offset`
	if err := ioutil.WriteFile(path, []byte(offsetStr), 0666); err != nil {
		f.log(`write offset error: ` + err.Error())
	}
	return offsetStr
}

func (f *File) readOffset() int64 {
	path := f.name + `/` + f.name + `.offset`
	if fs.NotExist(path) {
		return 0
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		f.log(fmt.Sprintf(`read offset %s: %s`, path, err.Error()))
		return 0
	}
	offset, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		f.log(fmt.Sprintf(`parse offset %s(%s): %s`, path, string(content), err.Error()))
		return 0
	}
	return offset
}

func (f *File) offset() int64 {
	if offset, err := f.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		f.log(`get offset error: ` + err.Error())
		return -1
	}
}

func (f *File) size() int64 {
	if fi, err := f.file.Stat(); err == nil {
		return fi.Size()
	} else {
		f.log(`get size error: ` + err.Error())
		return -1
	}
}
