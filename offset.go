package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/lovego/xiaomei/utils/fs"
)

type offsetInfo struct {
	FilePath string `json:"filepath"`
	Offset   int64  `json:"offset"`
}

const offsetDir = `/home/ubuntu/logs/logc/`
const offsetFile = `offset.json`

var offsetPath = path.Join(offsetDir, offsetFile)
var offsetMap = struct {
	m map[string]int64
	sync.RWMutex
}{m: make(map[string]int64)}

func init() {
	if err := os.MkdirAll(offsetDir, os.ModePerm); err != nil {
		panic(err)
	}
	if !fs.Exist(offsetPath) {
		_, err := os.Create(offsetPath)
		if err != nil {
			panic(err)
		}
	}
}

func getOffset(paths []string) map[string]int64 {
	data := readOffset()
	result := make(map[string]int64)
	for _, p := range paths {
		if data != nil {
			result[p] = data[p]
		} else {
			result[p] = 0
		}
	}
	return result
}

func updateOffset(infos []offsetInfo) bool {
	data := make(map[string]int64)
	for _, info := range infos {
		data[info.FilePath] = info.Offset
	}
	return writeOffset(data)
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
		printLog(`read offset error:`, printLog)
		return nil
	}
	return data
}

func writeOffset(data map[string]int64) bool {
	b, _ := json.Marshal(data)
	err := ioutil.WriteFile(offsetPath, b, 0666)
	if err != nil {
		printLog(`write offset error:`, err)
		return false
	}
	return true
}
