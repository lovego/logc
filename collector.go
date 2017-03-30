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

type File struct {
	Filepath string
	Offset   int64
	Org      string
	file     *os.File
}

func collector(paths []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					utils.Protect(func() {
						collect(event.Name)
					})
				}
			case err := <-watcher.Errors:
				writeLog("error:", err.Error())
			}
		}
	}()

	for _, filepath := range paths {
		err = watcher.Add(filepath)
		if err != nil {
			writeLog(filepath, `notify error:`, err.Error())
			continue
		}
		writeLog(`start notify file:`, filepath)
	}
	<-done
}

func collect(filepath string) {
	writeLog("change file:", filepath)
	monitorFiles.RLock()
	file := monitorFiles.data[filepath]
	monitorFiles.RUnlock()
	file.Collect()
}

func (f *File) Collect() {
	f.getFile()
	f.updateFiles()
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

func (f *File) read() string {
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

func (f *File) curOff(whence int) {
	off, err := f.file.Seek(0, whence)
	if err != nil {
		panic(err)
	}
	f.Offset = off
}

func (f *File) getFile() {
	if f.file != nil {
		return
	}
	file, err := os.Open(f.Filepath)
	if err != nil {
		panic(err)
	}
	f.file = file
}

func (f *File) updateFiles() {
	monitorFiles.RLock()
	monitorFiles.data[f.Filepath] = f
	monitorFiles.RUnlock()
}

func (f *File) checkFileOffset() {
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
