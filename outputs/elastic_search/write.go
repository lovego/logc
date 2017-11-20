package elastic_search

import (
	"strings"
	"time"

	loggerpkg "github.com/lovego/xiaomei/utils/logger"
	"github.com/nu7hatch/gouuid"
)

func (es *ElasticSearch) Write(
	rows []map[string]interface{}, logger *loggerpkg.Logger,
) (fatalError bool) {
	if len(rows) == 0 {
		return false
	}
	if es.timeSeriesIndex == nil {
		es.writeToIndex(es.index, rows)
		return false
	}

	return es.writeToTimeSeriesIndex(rows, logger)
}

func (es *ElasticSearch) writeToTimeSeriesIndex(
	rows []map[string]interface{}, logger *loggerpkg.Logger,
) (fatalError bool) {
	indicesRows, fatalError := es.timeSeriesIndex.Group(rows)
	if fatalError {
		return true
	}
	prune := false
	for _, one := range indicesRows {
		if one.Index != es.currentIndex {
			if es.ensureIndex(one.Index) {
				return true
			}
			es.currentIndex = one.Index
			prune = true
		}
		es.writeToIndex(one.Index, one.Rows)
	}
	if prune {
		es.timeSeriesIndex.Prune(es.client)
	}
	return false
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
