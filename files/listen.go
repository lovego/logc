package files

import (
	"encoding/json"
	"log"

	"github.com/lovego/xiaomei/utils"
	"gopkg.in/fsnotify.v1"
)

func (f *File) Listen() {
	log.Println(`listen ` + f.path)
	f.logger.Println(`listen ` + f.path)
	f.seekToSavedOffset()
	f.collect() // collect existing data before listen.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("notify new error: %v\n", err)
	}
	defer watcher.Close()

	if err := watcher.Add(f.path); err != nil {
		log.Printf("notify add %s error: %v\n", f.path, err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				utils.Protect(f.collect)
			}
		case err := <-watcher.Errors:
			f.logger.Printf("notify error: %v\n", err)
		}
	}
}

func (f *File) collect() {
	f.seekFrontIfTruncated()
	for rows := f.read(); len(rows) > 0; rows = f.read() {
		f.logd.Push(f.org, f.name, rows)
		offsetStr := f.writeOffset()
		f.logger.Printf("%d, %s\n", len(rows), offsetStr)
	}
}

func (f *File) read() []map[string]interface{} {
	rows := []map[string]interface{}{}
	for i := 0; i < 1000 && f.reader.More(); i++ {
		var row map[string]interface{}
		if err := f.reader.Decode(&row); err == nil {
			rows = append(rows, row)
		} else {
			f.logger.Println(`decode error: ` + err.Error())
			if _, ok := err.(*json.SyntaxError); ok {
				return rows
			}
		}
	}
	return rows
}
