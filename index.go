package main

import (
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/lovego/xiaomei/utils/httputil"
)

var monitorFiles = struct {
	sync.RWMutex
	data map[string]*File
}{data: make(map[string]*File)}

func FileInfo(orgName string) {
	fmt.Println(`remote address: `, remoteAddr)
	paths := []string{}
	httputil.Http(http.MethodGet, `http://`+path.Join(remoteAddr, `files?org=`+orgName), nil, nil, &paths)
	initFiles(orgName, paths)
	collector(paths)
}

func initFiles(orgName string, paths []string) {
	for _, filepath := range paths {
		monitorFiles.data[filepath] = &File{Filepath: filepath, Org: orgName}
	}
}
