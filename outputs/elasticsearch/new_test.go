package elasticsearch

import (
	"reflect"
	"testing"
)

func init() {
	theLogger = testLogger
}

func TestNew(t *testing.T) {
	for _, tc := range getNewTestCases() {
		got := New(`test.log`, tc.input, testLogger)
		expect := tc.expect
		if !reflect.DeepEqual(got, expect) {
			t.Fatalf("\ninput: %s\n expect: %+v\n    got: %+v\n", tc.input, tc.expect, got)
		}
	}
}

type testCaseT struct {
	input  map[string]interface{}
	expect *ElasticSearch
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

	return []testCaseT{
		{
			input: map[string]interface{}{
				"addrs":      []interface{}{"http://log-es.com/logc-test-"},
				"index":      "app-<2006.01.02>",
				"type":       "app-log",
				"mapping":    mappingIfc,
				"timeField":  "at",
				"timeFormat": "2006-01-02T15:04:05Z07:00",
				"indexKeep":  3,
			},
			expect: testES2,
		},
		{
			input: map[string]interface{}{
				"addrs":     []interface{}{"http://log-es.com/logc-test-"},
				"index":     "app-err",
				"type":      "app-err",
				"mapping":   mappingIfc,
				"timeField": "at",
			},
			expect: testES3,
		},
	}
}
