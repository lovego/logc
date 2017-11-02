package pusher

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lovego/logc/collector"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/lovego/xiaomei/utils/logger"
)

var httpClient = &httputil.Client{Client: http.DefaultClient}

type Getter struct {
	esIndex   string
	esType    string
	mergeJson string
}

func NewGetter(esAddrs []string, esIndex, esType, mergeJson string) collector.PusherGetter {
	if dataEs == nil {
		dataEs = elastic.New2(&httputil.Client{Client: http.DefaultClient}, esAddrs...)
	}
	return &Getter{esIndex, esType, mergeJson}
}

func (g *Getter) Get(log *logger.Logger) collector.Pusher {
	return &Pusher{Getter: g, logger: log}
}

type Pusher struct {
	*Getter
	logger *logger.Logger
}

func (p *Pusher) Push(rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	const max = 10 * time.Minute
	for interval := time.Second; ; {
		if p.push(rows) {
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

func (p *Pusher) push(docs []map[string]interface{}) bool {
	if p.mergeJson != `` {
		var merge map[string]interface{}
		if err := json.Unmarshal([]byte(p.mergeJson), &merge); err != nil {
			p.logger.Error("push data error: ", err)
			return false
		}
		for _, doc := range docs {
			for k, v := range merge {
				doc[k] = v
			}
		}
	}
	if err := dataEs.BulkCreate(p.esIndex+`/`+p.esType, convertDocs(docs)); err != nil {
		p.logger.Error("push data error: ", err)
		return false
	}
	return true
}

func convertDocs(docs []map[string]interface{}) [][2]interface{} {
	data := [][2]interface{}{}
	for _, doc := range docs {
		data = append(data, [2]interface{}{nil, doc})
	}
	return data
}
