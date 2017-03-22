package main

import (
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/httputil"
)

type fileInfo struct {
	FilePath string `json:"filepath"`
	Offset   int64  `json:"offset"`
	Fields   string `json:"fields"`
	SetId    string `json:"set_id"`
}

var monitorFiles = struct {
	sync.RWMutex
	data map[string]*File
}{data: make(map[string]*File)}

func FileInfo(token string) {
	fmt.Println(`remote address: `, remoteAddr)
	data := make(map[string]fileInfo)
	//token := `47a09e6c862b4daf4113b6ec52d70d6a`
	httputil.Http(http.MethodGet, `http://`+path.Join(remoteAddr, `files?token=`+token+`&ips=`+getIP()), nil, nil, &data)
	utils.PrintJson(data)
	initFiles(data)
	collector(data)
}

func initFiles(data map[string]fileInfo) {
	monitorFiles.Lock()
	for filepath, info := range data {
		monitorFiles.data[filepath] = &File{
			Filepath: filepath, Offset: info.Offset,
			SetId: info.SetId, Fields: info.Fields,
		}
	}
	monitorFiles.Unlock()
}
