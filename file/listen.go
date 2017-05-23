package file

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	pathpkg "path"
	"strconv"
	"time"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/httputil"
	"gopkg.in/fsnotify.v1"
)

func (f *File) Listen() {
	writeLog(`listen ` + f.path)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		writeLog(`notify new:`, err.Error())
	}
	defer watcher.Close()

	if err := watcher.Add(f.path); err != nil {
		writeLog(`notify add`, f.path+`:`, err.Error())
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				utils.Protect(func() {
					f.collect()
				})
			}
		case err := <-watcher.Errors:
			writeLog(`notify error:`, err.Error())
		}
	}
}

func (f *File) collect() {
	writeLog(`collect file:`, f.path)
	f.seekFrontIfTruncated()
	for rows := f.read(); len(rows) > 0; rows = f.read() {
		if f.push(rows) {
			writeLog(`the number of push rows:`, strconv.Itoa(len(rows)))
			if !f.writeOffset() {
				writeLog(f.path, `: update offset faild`)
			}
		} else {
			writeLog(`push faild`)
		}
	}
	writeLog(`collect complete`)
}

func (f *File) read() []map[string]interface{} {
	rows := []map[string]interface{}{}
	for i := 0; i < 1000 && f.reader.More(); i++ {
		var row map[string]interface{}
		if err := f.reader.Decode(&row); err != nil {
			writeLog(err.Error())
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

// 如果文件被截短，把文件seek到开头
func (f *File) seekFrontIfTruncated() {
	offset, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	info, err := os.Stat(f.path)
	if err != nil {
		panic(err)
	}
	if offset > info.Size() {
		f.file.Seek(0, os.SEEK_SET)
	}
}

func (f *File) push(rows []map[string]interface{}) bool {
	body := map[string]interface{}{`org`: f.org, `file`: f.name, `data`: rows}
	content, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	const max = time.Hour
	for interval := time.Second; ; {
		if push2Logd(content) {
			return true
		}
		time.Sleep(interval)
		if interval < max {
			interval *= 2
			if interval > max {
				interval = max
			}
		}
	}
	return false
}

var LogdAddr string

func push2Logd(content []byte) bool {
	result := make(map[string]string)
	uri := `http://` + pathpkg.Join(LogdAddr, `logs-data`)
	httputil.Http(http.MethodPost, uri, nil, bytes.NewBuffer(content), &result)
	return result[`msg`] == `ok`
}
