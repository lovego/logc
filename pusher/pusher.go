package pusher

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/lovego/xiaomei/utils/logger"
)

var httpClient = &httputil.Client{Client: http.DefaultClient}

type Getter struct {
	*config.File

	mergeJson string
}

func NewGetter(esAddrs []string, file *config.File, mergeJson string) collector.PusherGetter {
	if dataEs == nil {
		dataEs = elastic.New2(&httputil.Client{Client: http.DefaultClient}, esAddrs...)
	}
	return &Getter{file, mergeJson}
}

func (g *Getter) Get(log *logger.Logger) collector.Pusher {
	return &Pusher{Getter: g, logger: log}
}

type Pusher struct {
	*Getter

	logger *logger.Logger
}

var ensuredKeys = make(map[string]bool)

func (p *Pusher) Push(rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	const max = 10 * time.Minute
	data := p.groupRows(rows)
	for esIndex, d := range data {
		if !ensuredKeys[esIndex] {
			p.ensureIndex(esIndex)
			ensuredKeys[esIndex] = true
		}
		p.mergeJsonData(d)
		docs := convertDocs(d)
		for interval := time.Second; ; interval *= 2 {
			if docs = p.push(esIndex, docs); len(docs) == 0 {
				break
			}
			time.Sleep(interval)
			if interval > max {
				interval = max
			}
		}
	}
}

func (p *Pusher) push(esIndex string, docs [][2]interface{}) [][2]interface{} {
	if errs := dataEs.BulkCreate(esIndex+`/`+p.Type, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		p.logger.Error("push err is not elastic.BulkError type, but %T", errs)
	}
	return nil
}

func (p *Pusher) mergeJsonData(docs []map[string]interface{}) {
	if p.mergeJson != `` {
		var merge map[string]interface{}
		if err := json.Unmarshal([]byte(p.mergeJson), &merge); err != nil {
			p.logger.Error("merge json unexpected err: ", err)
			return
		}
		for _, doc := range docs {
			for k, v := range merge {
				doc[k] = v
			}
		}
	}
}

func (p *Pusher) groupRows(rows []map[string]interface{}) map[string][]map[string]interface{} {
	data := make(map[string][]map[string]interface{})
	for _, row := range rows {
		if value, ok := row[p.TimeField].(string); ok {
			if at, err := time.Parse(p.Parse, value); err == nil {
				esIndex := p.Index + `-` + at.Format(p.Layout)
				if data[esIndex] == nil {
					data[esIndex] = []map[string]interface{}{}
				}
				data[esIndex] = append(data[esIndex], row)
			} else {
				p.logger.Errorf("parse time field %s with layout %s error: %v", p.TimeField, p.Parse, err)
			}
		} else {
			p.logger.Errorf("time field %s not string", p.TimeField)
		}
	}
	return data
}

func convertDocs(docs []map[string]interface{}) [][2]interface{} {
	data := [][2]interface{}{}
	for _, doc := range docs {
		data = append(data, [2]interface{}{nil, doc})
	}
	return data
}
