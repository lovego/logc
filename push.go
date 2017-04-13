package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/lovego/xiaomei/utils/httputil"
)

const maxTime = time.Hour

func (f *file) push(data [][]string) bool {
	d := make(map[string]interface{})
	d[`org`] = f.org
	d[`filepath`] = f.path
	d[`data`] = data
	content, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	sleepTime := 1 * time.Second
	for success := pushRemote(content); !success; success = pushRemote(content) {
		time.Sleep(sleepTime)
		sleepTime *= 2
		if sleepTime >= maxTime {
			sleepTime = maxTime
		}
	}
	sleepTime = 1 * time.Second
	return true
}

func pushRemote(content []byte) bool {
	result := make(map[string]string)
	uri := `http://` + path.Join(remoteAddr, `logs-data`)
	httputil.Http(http.MethodPost, uri, nil, bytes.NewBuffer(content), &result)
	return result[`msg`] == `ok`
}
