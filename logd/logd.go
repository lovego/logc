package logd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/lovego/xiaomei/utils/httputil"
)

func (logd *Logd) Create(org string, files []map[string]interface{}) error {
	filesJson, err := json.Marshal(files)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Set(`org`, org)
	form.Set(`files`, string(filesJson))
	resp := struct {
		Code, Message string
	}{}
	if err := httputil.PostJson(logd.addr+`/org-files`, nil, form.Encode(), &resp); err != nil {
		return fmt.Errorf("create files error: %+v\n", err)
	}
	if resp.Code != `ok` {
		return fmt.Errorf("create files failed: %+v\n", resp)
	}
	return nil
}

func (logd *Logd) Push(org, file string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	pushUrl := logd.addr + `/logs-data?` + logd.getPushQuery(org, file)
	content, err := json.Marshal(rows)
	if err != nil {
		log.Printf("marshal rows error: %v\n", err)
		return
	}
	const max = time.Hour
	for interval := time.Second; ; {
		if push(pushUrl, content) {
			return
		}
		time.Sleep(interval)
		if interval < max {
			interval *= 2
			if interval > max {
				interval = max
			}
		}
	}
}

func (logd *Logd) getPushQuery(org, file string) string {
	query := url.Values{}
	query.Set(`org`, org)
	query.Set(`file`, file)
	if logd.mergeJson != `` {
		query.Set(`merge`, logd.mergeJson)
	}
	return query.Encode()
}

func push(pushUrl string, content []byte) bool {
	result := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}

	err := httputil.PostJson(pushUrl, nil, content, &result)
	if err != nil {
		log.Println("push data error: ", err)
		return false
	}
	if result.Code != `ok` {
		log.Printf("push data failed: %+v\n", result)
		return false
	}
	return true
}
