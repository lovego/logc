package elastic_search

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
)

// TODO: more error type retry
func (es *ElasticSearch) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := elasticSearch.BulkCreate(index+`/`+p.file.Type, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		logger.Error("push err is not elastic.BulkError type, but %T", errs)
	}
	return nil
}

func (es *ElasticSearch) ensureIndex(index string) (fatalError bool) {
	if err := elasticSearch.Ensure(index, nil); err != nil {
		logger.Errorf("ensure index %s error: %+v\n", index, err)
		return true
	}
	if err := elasticSearch.Put(index+`/_mapping/`+p.file.Type, map[string]interface{}{
		`properties`: p.file.Mapping,
	}, nil); err != nil {
		logger.Errorf("put mapping %s/%s error: %+v\n", index, p.file.Type, err)
		return true
	}
	return false
}
