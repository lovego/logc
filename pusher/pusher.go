package pusher

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/lovego/xiaomei/utils/httputil"
)

type Pusher struct {
	pushUrl string
	logger  *log.Logger
}

func New(addr, org, file, mergeJson string, logger *log.Logger) *Pusher {
	query := url.Values{}
	query.Set(`org`, org)
	query.Set(`file`, file)
	if mergeJson != `` {
		query.Set(`merge`, mergeJson)
	}
	return &Pusher{pushUrl: addr + `/logs-data?` + query.Encode(), logger: logger}
}

func (p *Pusher) Push(rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	content, err := json.Marshal(rows)
	if err != nil {
		p.logger.Printf("marshal rows error: %v\n", err)
		return
	}
	const max = time.Hour
	for interval := time.Second; ; {
		if p.push(content) {
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

func (p *Pusher) push(content []byte) bool {
	result := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}

	err := httputil.PostJson(p.pushUrl, nil, content, &result)
	if err != nil {
		p.logger.Println("push data error: ", err)
		return false
	}
	if result.Code != `ok` {
		p.logger.Printf("push data failed: %+v", result)
		return false
	}
	return true
}
