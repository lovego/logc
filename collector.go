package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/fsnotify.v1"
)

type File struct {
	Filepath string
	Fields   string
	Offset   int64
	SetId    string
	file     *os.File
	reader   *csv.Reader
}

func collector(infos map[string]fileInfo) {
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

	for filepath, _ := range infos {
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

	f.getReader()
	fields := parseFields(f.Fields)
	data := f.read(fields)
	for len(data) > 0 {
		printLog(`the number of push data:`, len(data))
		f.push(data)
		data = f.read(fields)
	}
	printLog(`all data has been pushed`)
}

func (f *File) read(fields [][2]string) []map[string]interface{} {
	data := []map[string]interface{}{}

	for i := 0; i < 1000; i++ {
		row, err := f.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			printLog(err)
			continue
		}

		if len(row) != len(fields) {
			continue
		}
		d := make(map[string]interface{})
		for j, fieldInfo := range fields {
			d[fieldInfo[0]] = row[j]
		}
		data = append(data, d)
	}
	f.curOff()
	return data
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

func (f *File) getReader() {
	if f.reader != nil {
		return
	}
	f.reader = csv.NewReader(f.file)
	f.reader.Comma = ' '
}

func (f *File) checkFileOffset() {
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
