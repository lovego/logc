package elastic_search

import (
	"strings"

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
			if !es.ensureIndex(one.Index, es.logger) {
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
	docs := es.addDocId(rows)
	var t Timer
	for {
		if docs = es.bulkCreate(index, docs); len(docs) == 0 {
			break
		}
		t.Sleep()
	}
}

func (es *ElasticSearch) addDocId(rows []map[string]interface{}) [][2]interface{} {
	docs := [][2]interface{}{}
	for _, doc := range rows {
		if uid, err := uuid.NewV4(); err != nil {
			es.logger.Errorf("generate uuid error: %v", err)
			docs = append(docs, [2]interface{}{nil, doc})
		} else {
			docs = append(docs, [2]interface{}{strings.Replace(uid.String(), `-`, ``, -1), doc})
		}
	}
	return docs
}
