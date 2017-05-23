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

func (f *File) writeOffset() bool {
	offset, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		writeLog(`get current offset:`, err.Error())
	}
	offsetStr := strconv.FormatInt(offset, 10)
	path := f.name + `/` + f.name + `.offset`
	if err := ioutil.WriteFile(path, []byte(offsetStr), 0666); err != nil {
		writeLog(`write offset error:`, err.Error())
		return false
	}
	return true
}
