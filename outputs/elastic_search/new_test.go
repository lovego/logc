package elastic_search

import (
	"os"
	"reflect"
	"testing"

	"github.com/lovego/logc/outputs/elastic_search/time_series_index"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

type testCaseT struct {
	input  map[string]interface{}
	expect *ElasticSearch
}

var testLogger = loggerpkg.New(``, os.Stderr, nil)

func TestNew(t *testing.T) {
	for _, tc := range getNewTestCases() {
		got := New(tc.input, `test.log`, testLogger)
		expect := tc.expect
		if expect != nil {
			expect.file = `test.log`
			expect.logger = loggerpkg.New(``, os.Stderr, nil)
			expect.client = got.client
		}
		if !reflect.DeepEqual(got, expect) {
			t.Fatalf("\ninput: %s\n expect: %+v\n    got: %+v\n", tc.input, tc.expect, got)
		}
	}
}

func getNewTestCases() []testCaseT {
	mappingIfc := map[interface{}]interface{}{
		"host":     map[interface{}]interface{}{"type": "keyword"},
		"query":    map[interface{}]interface{}{"type": "text"},
		"status":   map[interface{}]interface{}{"type": "keyword"},
		"req_body": map[interface{}]interface{}{"type": "integer"},
		"res_body": map[interface{}]interface{}{"type": "integer"},
		"agent":    map[interface{}]interface{}{"type": "text"},
		"at":       map[interface{}]interface{}{"type": "date"},
		"method":   map[interface{}]interface{}{"type": "keyword"},
		"path": map[interface{}]interface{}{
			"type": "text",
			"fields": map[interface{}]interface{}{
				"raw": map[interface{}]interface{}{"type": "keyword"},
			},
		},
		"ip":       map[interface{}]interface{}{"type": "ip"},
		"refer":    map[interface{}]interface{}{"type": "text"},
		"proto":    map[interface{}]interface{}{"type": "keyword"},
		"duration": map[interface{}]interface{}{"type": "float"},
	}

	mapping := map[string]map[string]interface{}{
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

	tsi, _ := time_series_index.New("app-<2006.01.02>", `at`, "2006-01-02T15:04:05Z0700", 3, testLogger)
	return []testCaseT{
		{
			input: map[string]interface{}{
				"addrs":      []interface{}{"http://log-es.com/logc-test-"},
				"index":      "app-<2006.01.02>",
				"type":       "app-log",
				"mapping":    mappingIfc,
				"timeField":  "at",
				"timeFormat": "2006-01-02T15:04:05Z0700",
				"indexKeep":  3,
			},
			expect: &ElasticSearch{
				addrs:           []string{`http://log-es.com/logc-test-`},
				index:           "app-<2006.01.02>",
				typ:             "app-log",
				mapping:         mapping,
				timeSeriesIndex: tsi,
			},
		},
		{
			input: map[string]interface{}{
				"addrs":     []interface{}{"http://log-es.com/logc-test-"},
				"index":     "app-err",
				"type":      "app-err",
				"mapping":   mappingIfc,
				"timeField": "at",
			},
			expect: &ElasticSearch{
				addrs:   []string{`http://log-es.com/logc-test-`},
				index:   "app-err",
				typ:     "app-err",
				mapping: mapping,
			},
		},
	}
}
