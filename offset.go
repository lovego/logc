package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lovego/xiaomei/utils/fs"
)

type offsetInfo struct {
	FilePath string `json:"filepath"`
	Offset   int64  `json:"offset"`
}

const offsetPath = `logc/offset.json`

var offsetData = struct {
	m map[string]int64
	sync.RWMutex
}{m: make(map[string]int64)}

func init() {
	if err := os.MkdirAll(filepath.Dir(offsetPath), os.ModePerm); err != nil {
		panic(err)
	}
	if !fs.Exist(offsetPath) {
		_, err := os.Create(offsetPath)
		if err != nil {
			panic(err)
		}
	}
	initOffset()
}

func initOffset() {
	m := readOffset()
	if m != nil {
		offsetData.m = m
	}
}

func updateOffset(path string, offset int64) bool {
	offsetData.Lock()
	offsetData.m[path] = offset
	success := writeOffset(offsetData.m)
	offsetData.Unlock()
	return success
}

func readOffset() map[string]int64 {
	b, err := ioutil.ReadFile(offsetPath)
	if err != nil {
		panic(err)
	}
	content := string(b)
	if strings.TrimSpace(content) == `` {
		return nil
	}
	data := make(map[string]int64)
	err = json.Unmarshal(b, &data)
	if err != nil {
		writeLog(`read offset error:`, err.Error())
		return nil
	}
	return data
}

func writeOffset(data map[string]int64) bool {
	b, _ := json.Marshal(data)
	err := ioutil.WriteFile(offsetPath, b, 0666)
	if err != nil {
		writeLog(`write offset error:`, err.Error())
		return false
	}
	return true
}
