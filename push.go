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

func (file *File) push(content string) bool {
	d := make(map[string]string)
	d[`org`] = file.Org
	d[`filepath`] = file.Filepath
	d[`content`] = content
	data, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	sleepTime := 1 * time.Second
	for success := pushRemote(data); !success; success = pushRemote(data) {
		time.Sleep(sleepTime)
		sleepTime *= 2
		if sleepTime > timeout {
			writeLog("collect faild.\n", content)
			sleepTime = 1 * time.Second
			return false
		}
	}
	sleepTime = 1 * time.Second
	return true
}

func pushRemote(data []byte) bool {
	result := make(map[string]string)
	uri := `http://` + path.Join(remoteAddr, `logs-data`)
	httputil.Http(http.MethodPost, uri, nil, bytes.NewBuffer(data), &result)
	return result[`msg`] == `ok`
}
