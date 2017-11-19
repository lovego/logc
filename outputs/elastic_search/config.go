package elastic_search

import (
	"errors"
	"github.com/spf13/cast"
)

type ElasticSearch struct {
	file            string
	addrs           []string
	index           string
	timeSeriesIndex *timeSeriesIndex
	indexKeep       int
	typ             string
	mapping         map[string]map[string]interface{}
	timeField       string
	timeFormat      string
}

func New(conf map[string]interface{}, file string) *ElasticSearch {
	if len(conf) == 0 {
		logger.Errorf(`elastic-search(%s): empty config.`, file)
		return nil
	}
	es := ElasticSearch{file: file}

	if addrs := parseAddrs(conf[`addrs`]); len(addrs) > 0 {
		es.addrs = addrs
	} else {
		return nil
	}

	var err error
	if file.TimeSeriesIndex, err = parseTimeSeriesIndex(file.Index); err != nil {
		log.Fatalf("file %s: %v", file.Path, err)
	}

	if file.Type == `` {
		log.Fatalf("type missing for file: %s", file.Path)
	}
	if file.TimeSeriesIndex != nil && file.TimeField == `` {
		log.Fatalf("timeField missing for file: %s, which has a time series index.", file.Path)
	}
	if file.TimeField != `` && file.TimeFormat == `` {
		file.TimeFormat = time.RFC3339
	}
	file.cleanMapping()
}

func parseAddrs(v interface{}) []string {
	if addrs, err := cast.ToStringSliceE(v); err == nil {
		if len(addrs) > 0 {
			return addrs
		} else {
			logger.Errorf(`elastic-search(%s): addrs is emtpty.`, file)
		}
	} else {
		logger.Errorf(`elastic-search(%s): addrs should be an string array.`, file)
	}
	return nil
}

func (file *File) cleanMapping() {
	for _, field := range file.Mapping {
		for k, v := range field {
			if data, ok := v.(map[interface{}]interface{}); ok {
				field[k] = convertMapKeyToStr(data)
			}
		}
	}
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
