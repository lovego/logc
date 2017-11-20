package elastic_search

import (
	"github.com/lovego/logc/outputs/elastic_search/time_series_index"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/spf13/cast"
)

type ElasticSearch struct {
	file    string
	addrs   []string
	index   string
	typ     string
	mapping map[string]map[string]interface{}
	client  *elastic.ES

	timeSeriesIndex *timeSeriesIndex
	currentIndex    string
}

func New(conf map[string]interface{}, file string) *ElasticSearch {
	if len(conf) == 0 {
		logger.Errorf(`elastic-search(%s): empty config.`, file)
		return nil
	}

	es := &ElasticSearch{file: file}

	var timeField, timeFormat string
	var indexKeep int
	if !es.parseConf(conf, &timeField, &timeFormat, &indexKeep) {
		return nil
	}
	if !es.check() {
		return nil
	}

	es.client = elastic.New2(&httputil.Client{Client: http.DefaultClient}, es.addrs...)

	if tsi, err := time_series_index.New(es.index, timeField, timeFormat, indexKeep); err == nil {
		es.timeSeriesIndex = tsi
	} else {
		logger.Errorf("elastic-search(%s) config: %v", es.file, err)
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
				logger.Errorf(`elastic-search(%s) config: addrs should be an string array.`, es.file)
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
				case `timeValue`:
					*timeValue = value
				}
			} else {
				logger.Errorf(`elastic-search(%s) config: %s should be a string.`, es.file, k)
				return false
			}
		case `mapping`:
			if mapping, ok := v.(map[interface{}]interface{}); ok {
				es.mapping = convertMapKeyToStr(mapping)
			}
			logger.Errorf(`elastic-search(%s) config: mapping should be a map.`, es.file)
			return false
		case `indexKeep`:
			if keep, ok := v.(int); ok {
				*indexKeep = keep
			}
			logger.Errorf(`elastic-search(%s): indexKeep should be an integer.`, es.file)
			return false
		default:
			logger.Errorf(`elastic-search(%s) config: unknown key: %s.`, es.file, k)
			return false
		}
	}
	return true
}

func (es *ElasticSearch) checkConf() bool {
	if len(es.addrs) == 0 {
		logger.Errorf(`elastic-search(%s) config: addrs is emtpty.`, es.file)
		return false
	}
	if es.index == `` {
		logger.Errorf(`elastic-search(%s) config: empty index.`, es.file)
		return false
	}
	if es.typ == `` {
		logger.Errorf(`elastic-search(%s): empty type.`, es.file)
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
