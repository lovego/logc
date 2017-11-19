package elastic_search

import (
	"strings"
	"time"

	"github.com/lovego/xiaomei/utils/logger"
	"github.com/nu7hatch/gouuid"
)

type Pusher struct {
	logger       *logger.Logger
	currentIndex string
}

type indexRows struct {
	index string
	rows  []map[string]interface{}
}

func (p *Pusher) Push(rows []map[string]interface{}) bool {
	if len(rows) == 0 {
		return true
	}
	if p.file.TimeSeriesIndex != nil {
		for _, indexData := range p.getTimeSeriesIndicesRows(rows) {
			if indexData.index != p.currentIndex {
				if !p.ensureIndex(indexData.index) {
					return false
				}
				p.currentIndex = indexData.index
			}
			p.push(indexData.index, indexData.rows)
		}
	} else {
		p.push(p.file.Index, rows)
	}
	return true
}

func (p *Pusher) push(index string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	p.mergeJsonData(rows)
	docs := p.convertDocs(rows)
	for {
		if docs = p.bulkCreate(index, docs); len(docs) == 0 {
			break
		}

		var interval time.Duration
		const max = 10 * time.Minute
		if interval <= 0 {
			interval = time.Second
		} else if interval < max {
			if interval *= 2; interval > max {
				interval = max
			}
		}
		time.Sleep(interval)
	}
}

func (p *Pusher) mergeJsonData(rows []map[string]interface{}) {
	if len(conf.MergeData) > 0 {
		for _, row := range rows {
			for k, v := range conf.MergeData {
				row[k] = v
			}
		}
	}
}

func (p *Pusher) getTimeSeriesIndicesRows(rows []map[string]interface{}) (result []indexRows) {
	indices := []string{}
	m := make(map[string][]map[string]interface{})
	for _, row := range rows {
		if index := p.getTimeSeriesIndexName(row); index != `` {
			if m[index] == nil {
				m[index] = []map[string]interface{}{row}
				indices = append(indices, index)
			} else {
				m[index] = append(m[index], row)
			}
		}
	}
	for _, index := range indices {
		result = append(result, indexRows{index: index, rows: m[index]})
	}
	return
}

func (p *Pusher) getTimeSeriesIndexName(row map[string]interface{}) string {
	value, ok := row[p.file.TimeField].(string)
	if !ok {
		p.logger.Errorf("non string timeField %s: %v", p.file.TimeField, row[p.file.TimeField])
		return ``
	}
	at, err := time.Parse(p.file.TimeFormat, value)
	if err != nil {
		p.logger.Errorf("parse timeField %s with layout %s error: %v", p.file.TimeField, p.file.TimeFormat, err)
		return ``
	}
	return p.file.TimeSeriesIndex.Get(at)
}

func (p *Pusher) convertDocs(docs []map[string]interface{}) [][2]interface{} {
	data := [][2]interface{}{}
	for _, doc := range docs {
		if id, err := genUUID(); err != nil {
			p.logger.Errorf("generate uuid error: %v", err)
			data = append(data, [2]interface{}{nil, doc})
		} else {
			data = append(data, [2]interface{}{id, doc})
		}
	}
	return data
}

func genUUID() (string, error) {
	if uid, err := uuid.NewV4(); err != nil {
		return ``, err
	} else {
		return strings.Replace(uid.String(), `-`, ``, -1), nil
	}
}