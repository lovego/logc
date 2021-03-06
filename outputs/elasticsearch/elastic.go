package elasticsearch

import (
	"github.com/lovego/elastic"
	loggerpkg "github.com/lovego/logger"
)

func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := es.client.BulkIndex(index+`/_doc`, docs); errs != nil {
		bulkErr, _ := errs.(elastic.BulkError)
		var failedItems [][2]interface{}
		if bulkErr != nil {
			failedItems = bulkErr.FailedItems(false, false)
		}
		es.logger.Errorf(
			"%s: bulkIndex error: %v\nfailed items: %v", es.collectorId, errs, failedItems,
		)
		if bulkErr != nil {
			return bulkErr.FailedItems(true, true)
		}
		return docs
	}
	return nil
}

func (es *ElasticSearch) ensureIndex(index string, logger *loggerpkg.Logger) bool {
	es.logger.Infof("check index: %s", index)
	if ok, err := es.client.Exist(index); err != nil {
		logger.Errorf("%s: check index %s existence error: %+v\n", es.collectorId, index, err)
		return false
	} else if !ok {
		es.logger.Infof("create index: %s", index)
		if err := es.client.Put(index, nil, nil); err != nil && !elastic.IsIndexAreadyExists(err) {
			logger.Errorf("%s: create index %s error: %+v\n", es.collectorId, index, err)
			return false
		}
	}
	// update mapping independently to allow new fields being added.
	if len(es.mapping) == 0 {
		return true
	}
	if err := es.client.Put(index+`/_mapping`, es.mapping, nil); err != nil {
		logger.Errorf("%s: put mapping for %s error: %+v\n", es.collectorId, index, err)
		return false
	}
	return true
}
