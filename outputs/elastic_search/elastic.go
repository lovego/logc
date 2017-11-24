package elastic_search

import (
	"github.com/lovego/elastic"
	loggerpkg "github.com/lovego/logger"
)

// TODO: more error type retry
func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := es.client.BulkIndex(index+`/`+es.typ, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		es.logger.Errorf("bulkCreate err: %v", errs)
		return docs
	}
	return nil
}

func (es *ElasticSearch) ensureIndex(index string, logger *loggerpkg.Logger) bool {
	es.logger.Printf("check index: %s", index)
	if ok, err := es.client.Exist(index); err != nil {
		logger.Errorf("check index %s existence error: %+v\n", index, err)
		return false
	} else if !ok {
		es.logger.Printf("create index: %s", index)
		if err := es.client.Put(index, nil, nil); err != nil {
			logger.Errorf("create index %s error: %+v\n", index, err)
			return false
		}
	}
	if err := es.client.Put(index+`/_mapping/`+es.typ, map[string]interface{}{
		`properties`: es.mapping,
	}, nil); err != nil {
		logger.Errorf("put mapping %s/%s error: %+v\n", index, es.typ, err)
		return false
	}
	return true
}
