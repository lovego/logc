package logd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/httputil"
)

func (logd *Logd) FilesOf(org string) (files map[string]string) {
	query := url.Values{}
	query.Set(`org`, org)
	filesUrl := logd.addr + `/files?` + query.Encode()
	httputil.Http(http.MethodGet, filesUrl, nil, nil, &files)
	return
}

func (logd *Logd) Push(org, file string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	pushUrl := logd.addr + `/logs-data?` + logd.getPushQuery(org, file)
	content, err := json.Marshal(rows)
	if err != nil {
		utils.Logf(`marshal rows error: %v`, err)
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
		Code string `json:"code"`
	}{}

	httputil.Http(http.MethodPost, pushUrl, nil, bytes.NewBuffer(content), &result)
	return result.Code == `ok`
}
