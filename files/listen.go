package files

import (
	"bytes"
	"encoding/json"
	"io"
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
	var line []byte
	for size := 0; size < 2*1024*1024; size += len(line) {
		var err error
		line, err = f.reader.ReadBytes('\n')
		if len(line) > 0 {
			var row map[string]interface{}
			if err := json.Unmarshal(line, &row); err == nil {
				rows = append(rows, row)
			} else {
				if line = bytes.TrimSpace(line); len(line) > 0 {
					f.logger.Printf("json error(%v): %s\n", err, line)
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				f.logger.Println(`read error: ` + err.Error())
			}
			return rows
		}
	}
	return rows
}
