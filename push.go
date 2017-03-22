package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/bughou-go/xiaomei/utils/httputil"
)

const timeout = time.Hour

func (file *File) push(data []map[string]interface{}) {
	d := make(map[string]interface{})
	d[`offset`] = strconv.FormatInt(file.Offset, 10)
	d[`set_id`] = file.SetId
	d[`ips`] = getIP()
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
			printLog("collect faild.", string(content))
			break
		}
	}
	sleepTime = 1 * time.Second
}

func pushRemote(content []byte) bool {
	uri := `http://` + path.Join(remoteAddr, `logs-data`)
	status, err := httputil.HttpStatus(http.MethodPost, uri, nil, bytes.NewBuffer(content))
	if err != nil {
		printLog(`push data error: `, err)
	}
	return status == 200
}
