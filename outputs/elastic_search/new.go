package elastic_search

import (
	"net/http"

	"github.com/lovego/logc/outputs/elastic_search/time_series_index"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
	"github.com/spf13/cast"
)

var theLogger *loggerpkg.Logger

func Setup(logger *loggerpkg.Logger) {
	theLogger = logger
}

type ElasticSearch struct {
	collectorId string
	addrs       []string
	index       string
	typ         string
	mapping     map[string]map[string]interface{}
	client      *elastic.ES

	timeSeriesIndex *time_series_index.TimeSeriesIndex
	currentIndex    string
	logger          *loggerpkg.Logger
}

func New(collectorId string, conf map[string]interface{}, logger *loggerpkg.Logger) *ElasticSearch {
	if len(conf) == 0 {
		theLogger.Errorf(`elastic-search(%s): empty config.`, collectorId)
		return nil
	}

	es := &ElasticSearch{collectorId: collectorId, logger: logger}

	var timeField, timeFormat string
	var indexKeep int
	if !es.parseConf(conf, &timeField, &timeFormat, &indexKeep) {
		return nil
	}
	if !es.checkConf() {
		return nil
	}

	if tsi, err := time_series_index.New(
		es.index, timeField, timeFormat, indexKeep, logger,
	); err == nil {
		es.timeSeriesIndex = tsi
	} else {
		theLogger.Errorf("elastic-search(%s) config: %v", es.collectorId, err)
		return nil
	}
	es.client = elastic.New2(&httputil.Client{Client: http.DefaultClient}, es.addrs...)
	if es.timeSeriesIndex == nil && !es.ensureIndex(es.index, theLogger) {
		return nil
	}

	return es
}

func (es *ElasticSearch) parseConf(conf map[string]interface{},
	timeField, timeFormat *string, indexKeep *int) bool {
	for k, v := range conf {
		switch k {
		case `addrs`:
			if addrs, err := cast.ToStringSliceE(v); err == nil {
				es.addrs = addrs
			} else {
				theLogger.Errorf(`elastic-search(%s) config: addrs should be an string array.`, es.collectorId)
				return false
			}
		case `index`, `type`, `timeField`, `timeFormat`:
			if value, ok := v.(string); ok {
				switch k {
				case `index`:
					es.index = value
				case `type`:
					es.typ = value
				case `timeField`:
					*timeField = value
				case `timeFormat`:
					*timeFormat = value
				}
			} else {
				theLogger.Errorf(`elastic-search(%s) config: %s should be a string.`, es.collectorId, k)
				return false
			}
		case `mapping`:
			if !es.parseMapping(v) {
				return false
			}
		case `indexKeep`:
			if keep, ok := v.(int); ok {
				*indexKeep = keep
			} else {
				theLogger.Errorf(`elastic-search(%s) config: indexKeep should be an integer.`, es.collectorId)
				return false
			}
		case `@collectorId`, `@type`:
		default:
			theLogger.Errorf(`elastic-search(%s) config: unknown key: %s.`, es.collectorId, k)
			return false
		}
	}
	return true
}

func (es *ElasticSearch) parseMapping(v interface{}) bool {
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		theLogger.Errorf(`elastic-search(%s) config: mapping should be a map.`, es.collectorId)
		return false
	}

	mapping := make(map[string]map[string]interface{})
	for k, v := range m {
		kk, ok := k.(string)
		if !ok {
			theLogger.Errorf(`elastic-search(%s) config: invalid mapping.`, es.collectorId)
			return false
		}
		vv, ok := v.(map[interface{}]interface{})
		if !ok {
			theLogger.Errorf(`elastic-search(%s) config: invalid mapping.`, es.collectorId)
			return false
		}
		mapping[kk] = convertMapKeyToStr(vv)
	}
	es.mapping = mapping
	return true
}

func (es *ElasticSearch) checkConf() bool {
	if len(es.addrs) == 0 {
		theLogger.Errorf(`elastic-search(%s) config: addrs is emtpty.`, es.collectorId)
		return false
	}
	if es.index == `` {
		theLogger.Errorf(`elastic-search(%s) config: index is empty.`, es.collectorId)
		return false
	}
	if es.typ == `` {
		theLogger.Errorf(`elastic-search(%s) config: type is empty.`, es.collectorId)
		return false
	}
	return true
}

func convertMapKeyToStr(data map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		if mapData, ok := v.(map[interface{}]interface{}); ok {
			result[k.(string)] = convertMapKeyToStr(mapData)
		} else {
			result[k.(string)] = v
		}
	}
	return result
}