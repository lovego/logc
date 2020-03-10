package elasticsearch

import (
	"strings"

	"log"

	"github.com/nu7hatch/gouuid"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// return false to stop collector.
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
			var t Timer
			for !es.ensureIndex(one.Index, es.logger) {
				t.Sleep()
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
	if es.addTypeSuffix {
		arrayAddTypeSuffixToMapKeys(rows, es.mapping)
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
		convertKeyWithDot(doc)
		if uid, err := uuid.NewV4(); err != nil {
			es.logger.Errorf("generate uuid error: %v", err)
			docs = append(docs, [2]interface{}{nil, doc})
		} else {
			docs = append(docs, [2]interface{}{strings.Replace(uid.String(), `-`, ``, -1), doc})
		}
	}
	return docs
}

// convert dot(.) in key to underline(_)
func convertKeyWithDot(doc map[string]interface{}) {
	for key, value := range doc {
		if strings.ContainsRune(key, '.') {
			newKey := strings.Replace(key, `.`, `_`, -1)
			doc[newKey] = value
			delete(doc, key)
		}
		if v, ok := value.(map[string]interface{}); ok {
			convertKeyWithDot(v)
		}
		if vs, ok := value.([]map[string]interface{}); ok {
			for _, v := range vs {
				convertKeyWithDot(v)
			}
		}
	}
}
