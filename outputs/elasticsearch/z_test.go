package elasticsearch

import (
	"net/http"
	"os"
	"time"

	"github.com/lovego/elastic"
	"github.com/lovego/httputil"
	"github.com/lovego/logc/outputs/elasticsearch/time_series_index"
	loggerpkg "github.com/lovego/logger"
)

var testEsAddrs = []string{`http://127.0.0.1:9200/logc-test-`}
var testEsClient = elastic.New2(&httputil.Client{Client: http.DefaultClient}, testEsAddrs...)

var testMapping = map[string]interface{}{
	"properties": map[string]interface{}{
		"host":          map[string]interface{}{"type": "keyword"},
		"status":        map[string]interface{}{"type": "keyword"},
		"req_body_size": map[string]interface{}{"type": "integer"},
		"res_body_size": map[string]interface{}{"type": "integer"},
		"agent":         map[string]interface{}{"type": "text"},
		"at":            map[string]interface{}{"type": "date"},
		"method":        map[string]interface{}{"type": "keyword"},
		"path": map[string]interface{}{
			"type": "text",
			"fields": map[string]interface{}{
				"raw": map[string]interface{}{"type": "keyword"},
			},
		},
		"ip":       map[string]interface{}{"type": "ip"},
		"refer":    map[string]interface{}{"type": "text"},
		"proto":    map[string]interface{}{"type": "keyword"},
		"duration": map[string]interface{}{"type": "float"},
	},
	"dynamic_templates": []interface{}{
		map[string]interface{}{
			"query": map[string]interface{}{
				"path_match": "query.*",
				"mapping": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{"type": "keyword"},
					},
				},
			},
		},
	},
}

var testLogger = loggerpkg.New(``, os.Stderr, nil)

var testTsi1, _ = time_series_index.New("", "app-<2006.01.02>", `at`, time.RFC3339, 0, testLogger)
var testTsi2, _ = time_series_index.New("test.log", "app-<2006.01.02>", `at`, time.RFC3339, 3, testLogger)

var testES1 = &ElasticSearch{
	collectorId:     `test.log`,
	addrs:           testEsAddrs,
	index:           "app-<2006.01.02>",
	mapping:         testMapping,
	timeSeriesIndex: testTsi1,
	client:          testEsClient,
	logger:          testLogger,
}

var testES2 = &ElasticSearch{
	collectorId:     `test.log`,
	addrs:           testEsAddrs,
	index:           "app-<2006.01.02>",
	mapping:         testMapping,
	timeSeriesIndex: testTsi2,
	client:          testEsClient,
	logger:          testLogger,
}

var testES3 = &ElasticSearch{
	collectorId: `test.log`,
	addrs:       testEsAddrs,
	index:       "app-err",
	mapping:     testMapping,
	client:      testEsClient,
	logger:      testLogger,
}
