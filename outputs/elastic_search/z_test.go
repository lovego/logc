package elastic_search

import (
	"net/http"
	"os"
	"time"

	"github.com/lovego/logc/outputs/elastic_search/time_series_index"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

var testEsAddrs = []string{`http://log-es.com/logc-test-`}
var testEsClient = elastic.New2(&httputil.Client{Client: http.DefaultClient}, testEsAddrs...)

var testMapping = map[string]map[string]interface{}{
	"host":     {"type": "keyword"},
	"query":    {"type": "text"},
	"status":   {"type": "keyword"},
	"req_body": {"type": "integer"},
	"res_body": {"type": "integer"},
	"agent":    {"type": "text"},
	"at":       {"type": "date"},
	"method":   {"type": "keyword"},
	"path": {
		"type": "text",
		"fields": map[string]interface{}{
			"raw": map[string]interface{}{"type": "keyword"},
		},
	},
	"ip":       {"type": "ip"},
	"refer":    {"type": "text"},
	"proto":    {"type": "keyword"},
	"duration": {"type": "float"},
}

var testLogger = loggerpkg.New(``, os.Stderr, nil)

var testTsi1, _ = time_series_index.New("app-<2006.01.02>", `at`, time.RFC3339, 0, testLogger)
var testTsi2, _ = time_series_index.New("app-<2006.01.02>", `at`, time.RFC3339, 3, testLogger)

var testES1 = &ElasticSearch{
	file:            `test.log`,
	addrs:           testEsAddrs,
	index:           "app-<2006.01.02>",
	typ:             "app-log",
	mapping:         testMapping,
	timeSeriesIndex: testTsi1,
	client:          testEsClient,
	logger:          testLogger,
}

var testES2 = &ElasticSearch{
	file:            `test.log`,
	addrs:           testEsAddrs,
	index:           "app-<2006.01.02>",
	typ:             "app-log",
	mapping:         testMapping,
	timeSeriesIndex: testTsi2,
	client:          testEsClient,
	logger:          testLogger,
}

var testES3 = &ElasticSearch{
	file:    `test.log`,
	addrs:   testEsAddrs,
	index:   "app-err",
	typ:     "app-err",
	mapping: testMapping,
	client:  testEsClient,
	logger:  testLogger,
}
