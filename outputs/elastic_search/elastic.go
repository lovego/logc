package elastic_search

import (
	"net/http"
	"time"

	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
)

// TODO: more error type retry
func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := es.client.BulkCreate(index+`/`+es.typ, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		es.logger.Errorf("bulkCreate err: %v", errs)
		return docs
	}
	return nil
}

func (es *ElasticSearch) setupIndex() bool {
	es.client = elastic.New2(&httputil.Client{Client: http.DefaultClient}, es.addrs...)
	if es.timeSeriesIndex == nil {
		return es.ensureIndex(es.index)
	} else {
		return es.ensureIndex(es.timeSeriesIndex.Get(time.Now()))
	}
}

func (es *ElasticSearch) ensureIndex(index string) bool {
	if err := es.client.Ensure(index, nil); err != nil {
		es.logger.Errorf("ensure index %s error: %+v\n", index, err)
		return false
	}
	if err := es.client.Put(index+`/_mapping/`+es.typ, map[string]interface{}{
		`properties`: es.mapping,
	}, nil); err != nil {
		es.logger.Errorf("put mapping %s/%s error: %+v\n", index, es.typ, err)
		return false
	}
	return true
}
