package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

// 如果文件被截短，把文件seek到开头
func (f *File) seekFrontIfTruncated() {
	offset, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		f.log(`get current offset error: ` + err.Error())
		return
	}
	fi, err := f.file.Stat()
	if err != nil {
		f.log(`stat file error: ` + err.Error())
		return
	}
	if offset > fi.Size() {
		f.file.Seek(0, os.SEEK_SET)
	}
}

func (f *File) readOffset() int64 {
	path := f.name + `/` + f.name + `.offset`
	if !fs.Exist(path) {
		return 0
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		utils.Log(fmt.Sprintf(`read offset %s: %s`, path, err.Error()))
		return 0
	}
	offset, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		utils.Log(fmt.Sprintf(`parse offset %s(%s): %s`, path, string(content), err.Error()))
		return 0
	}
	return offset
}

func (f *File) writeOffset() string {
	offset, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		f.log(`get current offset error: ` + err.Error())
		return ``
	}
	offsetStr := strconv.FormatInt(offset, 10)

	path := f.name + `/` + f.name + `.offset`
	if err := ioutil.WriteFile(path, []byte(offsetStr), 0666); err != nil {
		f.log(`write offset error: ` + err.Error())
	}
	return offsetStr
}
