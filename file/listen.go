package file

import (
	"bytes"
	"encoding/json"
	"net/http"
	pathpkg "path"
	"strconv"
	"time"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/httputil"
	"gopkg.in/fsnotify.v1"
)

func (f *File) Listen() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Log(`notify new: ` + err.Error())
	}
	defer watcher.Close()

	if err := watcher.Add(f.path); err != nil {
		utils.Log(`notify add ` + f.path + `: ` + err.Error())
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				utils.Protect(f.Collect)
			}
		case err := <-watcher.Errors:
			f.Log(`notify error: ` + err.Error())
		}
	}
}

func (f *File) Collect() {
	f.seekFrontIfTruncated()
	for rows := f.read(); len(rows) > 0; rows = f.read() {
		f.push(rows)
		offsetStr := f.writeOffset()
		f.Log(strconv.Itoa(len(rows)) + `, ` + offsetStr)
	}
}

func (f *File) read() []map[string]interface{} {
	rows := []map[string]interface{}{}
	for i := 0; i < 1000 && f.reader.More(); i++ {
		var row map[string]interface{}
		if err := f.reader.Decode(&row); err == nil {
			rows = append(rows, row)
		} else {
			f.Log(`decode error: ` + err.Error())
		}
	}
	return rows
}

func (f *File) push(rows []map[string]interface{}) {
	body := map[string]interface{}{`org`: f.org, `file`: f.name, `data`: rows}
	content, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	const max = time.Hour
	for interval := time.Second; ; {
		if push2Logd(content) {
			return
		}
		time.Sleep(interval)
		if interval < max {
			interval *= 2
			if interval > max {
				interval = max
			}
		}
	}
}

var LogdAddr string

func push2Logd(content []byte) bool {
	result := make(map[string]string)
	uri := `http://` + pathpkg.Join(LogdAddr, `logs-data`)
	httputil.Http(http.MethodPost, uri, nil, bytes.NewBuffer(content), &result)
	return result[`msg`] == `ok`
}
