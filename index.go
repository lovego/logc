package main

import (
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/bughou-go/xiaomei/utils/httputil"
)

var monitorFiles = struct {
	sync.RWMutex
	data map[string]*File
}{data: make(map[string]*File)}

func FileInfo(orgName string) {
	fmt.Println(`remote address: `, remoteAddr)
	data := []string{}
	httputil.Http(http.MethodGet, `http://`+path.Join(remoteAddr, `files?org=`+orgName), nil, nil, &data)
	initFiles(orgName, data)
	collector(data)
}

func initFiles(orgName string, data []string) {
	monitorFiles.RLock()
	for _, filepath := range data {
		monitorFiles.data[filepath] = &File{Filepath: filepath, Org: orgName}
	}
	monitorFiles.RUnlock()
}
