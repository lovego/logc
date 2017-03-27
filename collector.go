package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

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
				printLog("error:", err)
			}
		}
	}()

	for _, filepath := range paths {
		err = watcher.Add(filepath)
		if err != nil {
			printLog(filepath, ` notify error:`, err)
			continue
		}
		printLog(`start notify file: `, filepath)
	}
	<-done
}

func collect(filepath string) {
	printLog("change file: ", filepath)
	monitorFiles.RLock()
	file := monitorFiles.data[filepath]
	monitorFiles.RUnlock()
	printLog(`start collect log`)
	file.Collect()
}

func (f *File) Collect() {
	f.getFile()
	f.checkFileOffset()
	f.file.Seek(f.Offset, os.SEEK_CUR)
	for content := f.read(); content != ``; content = f.read() {
		if f.push(content) {
			offsetData.RLock()
			offsetData.m[f.Filepath] = f.Offset
			offsetData.RUnlock()
			if !updateOffset() {
				printLog(f.Filepath, `: update offset faild`)
			}
		}
	}
	printLog(`all data has been pushed`)
	f.updateFiles()
}

func (f *File) read() string {
	var content string
	for {
		b := make([]byte, 1024*100)
		_, err := f.file.ReadAt(b, f.Offset)
		if err == io.EOF {
			break
		}
		if err != nil {
			printLog(err)
			continue
		}
		f.curOff()
		content += string(b)
	}
	return content
}

func (f *File) curOff() {
	off, err := f.file.Seek(0, os.SEEK_CUR)
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
	printLog(f.Filepath, " end offset: ", ret)
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
