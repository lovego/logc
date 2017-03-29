package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/lovego/xiaomei/utils/httputil"
)

const timeout = time.Hour

func (file *File) push(data string) bool {
	d := make(map[string]interface{})
	d[`org`] = file.Org
	d[`filepath`] = file.Filepath
	d[`data`] = data
	content, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	sleepTime := 1 * time.Second
	for success := pushRemote(content); !success; success = pushRemote(content) {
		time.Sleep(sleepTime)
		sleepTime *= 2
		if sleepTime > timeout {
			writeLog("collect faild.\n", data)
			sleepTime = 1 * time.Second
			return false
		}
	}
	sleepTime = 1 * time.Second
	return true
}

func pushRemote(content []byte) bool {
	data := make(map[string]string)
	uri := `http://` + path.Join(remoteAddr, `logs-data`)
	httputil.Http(http.MethodPost, uri, nil, bytes.NewBuffer(content), &data)
	return data[`msg`] == `ok`
}
