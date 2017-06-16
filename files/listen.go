package files

import (
	"github.com/lovego/xiaomei/utils"
	"gopkg.in/fsnotify.v1"
)

func (f *File) Listen() {
	f.seekToSavedOffset()
	f.collect() // collect existing data before listen.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logf("notify new error: %v", err)
	}
	defer watcher.Close()

	if err := watcher.Add(f.path); err != nil {
		utils.Logf("notify add %s error: %v", f.path, err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				utils.Protect(f.collect)
			}
		case err := <-watcher.Errors:
			f.Log(`notify error: %v`, err)
		}
	}
}

func (f *File) collect() {
	f.seekFrontIfTruncated()
	for rows := f.read(); len(rows) > 0; rows = f.read() {
		f.logd.Push(f.org, f.name, rows)
		offsetStr := f.writeOffset()
		f.Log("%d, %v", len(rows), offsetStr)
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
