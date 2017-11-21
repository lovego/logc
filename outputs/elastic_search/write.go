package elastic_search

import (
	"strings"
	"time"

	"github.com/nu7hatch/gouuid"
)

func (es *ElasticSearch) Write(rows []map[string]interface{}) bool {
	if len(rows) == 0 {
		return true
	}
	if es.timeSeriesIndex == nil {
		es.writeToIndex(es.index, rows)
		return true
	}

	return es.writeToTimeSeriesIndex(rows)
}

func (es *ElasticSearch) writeToTimeSeriesIndex(rows []map[string]interface{}) bool {
	indicesRows := es.timeSeriesIndex.Group(rows)
	if len(indicesRows) <= 0 { // Group encountered error.
		return false
	}
	prune := false
	for _, one := range indicesRows {
		if one.Index != es.currentIndex {
			if !es.ensureIndex(one.Index) {
				return false
			}
			es.currentIndex = one.Index
			prune = true
		}
		es.writeToIndex(one.Index, one.Rows)
	}
	if prune {
		es.timeSeriesIndex.Prune(es.client)
	}
	return true
}

func (es *ElasticSearch) writeToIndex(index string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	docs := es.convertDocs(rows)
	for {
		if docs = es.bulkCreate(index, docs); len(docs) == 0 {
			break
		}
		var t Timer
		t.Sleep()
	}
}

func (es *ElasticSearch) convertDocs(docs []map[string]interface{}) [][2]interface{} {
	data := [][2]interface{}{}
	for _, doc := range docs {
		if id, err := genUUID(); err != nil {
			es.logger.Errorf("generate uuid error: %v", err)
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

type Timer struct {
	duration time.Duration
}

func (t *Timer) Sleep() {
	const max = 10 * time.Minute
	if t.duration <= 0 {
		t.duration = time.Second
	} else if t.duration < max {
		if t.duration *= 2; t.duration > max {
			t.duration = max
		}
	}
	time.Sleep(t.duration)
}
