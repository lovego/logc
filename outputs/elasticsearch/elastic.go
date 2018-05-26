package elasticsearch

import (
	"github.com/lovego/elastic"
	loggerpkg "github.com/lovego/logger"
)

// TODO: more error type retry
func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := es.client.BulkIndex(index+`/_doc`, docs); errs != nil {
		es.logger.Errorf("%s: bulkIndex error: %v", es.collectorId, errs)
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		return docs
	}
	return nil
}

func (es *ElasticSearch) ensureIndex(index string, logger *loggerpkg.Logger) bool {
	es.logger.Printf("check index: %s", index)
	if ok, err := es.client.Exist(index); err != nil {
		logger.Errorf("%s: check index %s existence error: %+v\n", es.collectorId, index, err)
		return false
	} else if !ok {
		es.logger.Printf("create index: %s", index)
		if err := es.client.Put(index, nil, nil); err != nil && !elastic.IsIndexAreadyExists(err) {
			logger.Errorf("%s: create index %s error: %+v\n", es.collectorId, index, err)
			return false
		}
	}
	if err := es.client.Put(index+`/_mapping/_doc`, es.mapping, nil); err != nil {
		logger.Errorf("%s: put mapping %s/_doc error: %+v\n", es.collectorId, index, err)
		return false
	}
	return true
}
