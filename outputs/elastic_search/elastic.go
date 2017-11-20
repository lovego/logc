package elastic_search

import (
	"github.com/lovego/xiaomei/utils/elastic"
)

// TODO: more error type retry
func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := es.client.BulkCreate(index+`/`+es.typ, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		es.logger.Error("push err is not elastic.BulkError type, but %T", errs)
	}
	return nil
}

func (es *ElasticSearch) ensureIndex(index string) (fatalError bool) {
	if err := es.client.Ensure(index, nil); err != nil {
		es.logger.Errorf("ensure index %s error: %+v\n", index, err)
		return true
	}
	if err := es.client.Put(index+`/_mapping/`+es.typ, map[string]interface{}{
		`properties`: es.mapping,
	}, nil); err != nil {
		es.logger.Errorf("put mapping %s/%s error: %+v\n", index, es.typ, err)
		return true
	}
	return false
}
