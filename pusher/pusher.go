package pusher

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/lovego/logc/collector"
	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/lovego/xiaomei/utils/logger"
)

var httpClient = &httputil.Client{Client: http.DefaultClient}

type Getter struct {
	pushUrl string
}

func NewGetter(addr, org, file, mergeJson string) collector.PusherGetter {
	query := url.Values{}
	query.Set(`org`, org)
	query.Set(`file`, file)
	if mergeJson != `` {
		query.Set(`merge`, mergeJson)
	}
	return &Getter{pushUrl: addr + `/logs-data?` + query.Encode()}
}

func (g *Getter) Get(log *logger.Logger) collector.Pusher {
	return &Pusher{pushUrl: g.pushUrl, logger: log}
}

type Pusher struct {
	pushUrl string
	logger  *logger.Logger
}

func (p *Pusher) Push(rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	content, err := json.Marshal(rows)
	if err != nil {
		p.logger.Errorf("marshal rows error: %v\n", err)
		return
	}
	const max = 10 * time.Minute
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

	err := httpClient.PostJson(p.pushUrl, nil, content, &result)
	if err != nil {
		p.logger.Error("push data error: ", err)
		return false
	}
	if result.Code != `ok` {
		p.logger.Errorf("push data failed: %+v", result)
		return false
	}
	return true
}
