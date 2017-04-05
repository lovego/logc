package main

import (
	"bytes"
	"encoding/json"
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

// 如果文件被清空
func (f *file) collect() {
	writeLog(`collect file:`, f.path)
	f.checkFileOffset()
	f.file.Seek(f.Offset, os.SEEK_CUR)
	for content := f.read(); content != ``; content = f.read() {
		if f.push(content) {
			offsetData.RLock()
			offsetData.m[f.Filepath] = f.Offset
			offsetData.RUnlock()
			if !updateOffset() {
				writeLog(f.Filepath, `: update offset faild`)
			}
		} else {
			writeLog(`push faild`)
		}
	}
	writeLog(`collect complete`)
}

func (f *file) read() string {
	var content string
	for done := false; !done; {
		b := make([]byte, 1024*100)
		n, err := f.file.ReadAt(b, f.Offset)
		if err == io.EOF {
			done = true
		}
		if err != io.EOF && err != nil {
			writeLog(err.Error())
			continue
		}
		content += string(b[:n])
		if done {
			f.curOff(os.SEEK_END)
		} else {
			f.curOff(os.SEEK_CUR)
		}
	}
	return content
}

func (f *file) curOff(whence int) {
	off, err := f.file.Seek(0, whence)
	if err != nil {
		panic(err)
	}
	f.Offset = off
}

func (f *file) checkFileOffset() {
	offsetData.RLock()
	f.Offset = offsetData.m[f.Filepath]
	offsetData.RUnlock()
	ret, err := f.file.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	writeLog(f.Filepath, "end offset:", strconv.FormatInt(ret, 10))
	f.file.Seek(0, os.SEEK_SET)
	if ret < f.Offset {
		f.Offset = 0
	}
}

func parseFields(fieldsStr string) [][2]string {
	fields := [][2]string{}
	b := bytes.NewBufferString(fieldsStr).Bytes()
	err := json.Unmarshal(b, &fields)
	if err != nil {
		panic(err)
	}
	return fields
}
