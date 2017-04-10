package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/lovego/xiaomei/utils"
	"gopkg.in/fsnotify.v1"
)

type file struct {
	org    string
	path   string
	file   *os.File
	offset int64
	csv    *csv.Reader
}

func newFile(org, path string) *file {
	f := &file{
		org: org, path: path,
	}
	file, err := os.Open(path)
	if err != nil {
		writeLog(`open`, path+`:`, err.Error())
	}
	f.file = file
	f.offset = offsetData.m[path]
	return f
}

func (f *file) listen() {
	writeLog(`listen`, f.path)

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

func (f *file) collect() {
	writeLog(`collect file:`, f.path)
	f.csvReader()
	f.seekFrontIfTruncated()
	f.file.Seek(f.offset, os.SEEK_CUR)
	for data := f.read(); len(data) > 0; data = f.read() {
		if f.push(data) {
			writeLog(`the number of push data:`, strconv.Itoa(len(data)))
			offsetData.RLock()
			offsetData.m[f.path] = f.offset
			offsetData.RUnlock()
			if !updateOffset() {
				writeLog(f.path, `: update offset faild`)
			}
		} else {
			writeLog(`push faild`)
		}
	}
	writeLog(`collect complete`)
}

func (f *file) read() [][]string {
	data := [][]string{}
	for i := 0; i < 1000; i++ {
		row, err := f.csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeLog(err.Error())
			continue
		}
		data = append(data, row)
	}
	f.curOff()
	return data
}

func (f *file) curOff() {
	off, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	f.offset = off
}

func (f *file) csvReader() {
	if f.csv == nil {
		f.csv = csv.NewReader(f.file)
		f.csv.Comma = ' '
	}
}

// 如果文件被截短，把文件offset移动到开头
func (f *file) seekFrontIfTruncated() {
	ret, err := f.file.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	writeLog(f.path, "end offset:", strconv.FormatInt(ret, 10))
	f.file.Seek(0, os.SEEK_SET)
	if ret < f.offset {
		f.offset = 0
	}
}
