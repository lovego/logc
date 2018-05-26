package elasticsearch

import (
	"net/http"

	"github.com/lovego/elastic"
	"github.com/lovego/httputil"
	"github.com/lovego/logc/outputs/elasticsearch/time_series_index"
	loggerpkg "github.com/lovego/logger"
	"github.com/lovego/strmap"
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
	mapping     map[string]interface{}
	client      *elastic.ES

	timeSeriesIndex *time_series_index.TimeSeriesIndex
	currentIndex    string
	logger          *loggerpkg.Logger
}

func New(collectorId string, conf map[string]interface{}, logger *loggerpkg.Logger) *ElasticSearch {
	if len(conf) == 0 {
		theLogger.Errorf("elasticsearch(%s): empty config.", collectorId)
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
		collectorId, es.index, timeField, timeFormat, indexKeep, logger,
	); err == nil {
		es.timeSeriesIndex = tsi
	} else {
		theLogger.Errorf("elasticsearch(%s) config: %v", es.collectorId, err)
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
				theLogger.Errorf("elasticsearch(%s) config: addrs should be an string array.", es.collectorId)
				return false
			}
		case `index`, `type`, `timeField`, `timeFormat`:
			if value, ok := v.(string); ok {
				switch k {
				case `index`:
					es.index = value
				case `timeField`:
					*timeField = value
				case `timeFormat`:
					*timeFormat = value
				}
			} else {
				theLogger.Errorf("elasticsearch(%s) config: %s should be a string.", es.collectorId, k)
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
				theLogger.Errorf("elasticsearch(%s) config: indexKeep should be an integer.", es.collectorId)
				return false
			}
		case `@collectorId`, `@type`:
		default:
			theLogger.Errorf("elasticsearch(%s) config: unknown key: %s.", es.collectorId, k)
			return false
		}
	}
	return true
}

func (es *ElasticSearch) parseMapping(v interface{}) bool {
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		theLogger.Errorf("elasticsearch(%s) config: mapping should be a map.", es.collectorId)
		return false
	}
	if mapping, err := strmap.Convert(m, ""); err != nil {
		theLogger.Errorf("elasticsearch(%s) config: invalid mapping: %v.", es.collectorId, err)
		return false
	} else {
		es.mapping = mapping
	}
	return true
}

func (es *ElasticSearch) checkConf() bool {
	if len(es.addrs) == 0 {
		theLogger.Errorf("elasticsearch(%s) config: addrs is emtpty.", es.collectorId)
		return false
	}
	if es.index == `` {
		theLogger.Errorf("elasticsearch(%s) config: index is empty.", es.collectorId)
		return false
	}
	return true
}
