package elasticsearch

import (
	"strings"
	"testing"

	"github.com/lovego/deep"
)

func init() {
	theLogger = testLogger
	deep.CompareUnexportedFields = true
}

func TestNew(t *testing.T) {
	for _, tc := range getNewTestCases() {
		got := New(`test.log`, tc.input, testLogger)
		expect := tc.expect
		if diff := deep.Equal(got, expect); diff != nil {
			t.Fatalf("\ninput: %v\n diff: \n%+v\n", tc.input, strings.Join(diff, "\n"))
		}
	}
}

type testCaseT struct {
	input  map[string]interface{}
	expect *ElasticSearch
}

func getNewTestCases() []testCaseT {
	mappingIfc := map[interface{}]interface{}{
		"properties": map[interface{}]interface{}{
			"host":          map[interface{}]interface{}{"type": "keyword"},
			"status":        map[interface{}]interface{}{"type": "keyword"},
			"req_body_size": map[interface{}]interface{}{"type": "integer"},
			"res_body_size": map[interface{}]interface{}{"type": "integer"},
			"agent":         map[interface{}]interface{}{"type": "text"},
			"at":            map[interface{}]interface{}{"type": "date"},
			"method":        map[interface{}]interface{}{"type": "keyword"},
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
		},
		"dynamic_templates": []interface{}{
			map[interface{}]interface{}{
				"query": map[interface{}]interface{}{
					"path_match": "query.*",
					"mapping": map[interface{}]interface{}{
						"type": "text",
						"fields": map[interface{}]interface{}{
							"raw": map[interface{}]interface{}{"type": "keyword"},
						},
					},
				},
			},
		},
	}

	return []testCaseT{
		{
			input: map[string]interface{}{
				"addrs":      []interface{}{"http://127.0.0.1:9200/logc-test-"},
				"index":      "app-<2006.01.02>",
				"mapping":    mappingIfc,
				"timeField":  "at",
				"timeFormat": "2006-01-02T15:04:05Z07:00",
				"indexKeep":  3,
			},
			expect: testES2,
		},
		{
			input: map[string]interface{}{
				"addrs":     []interface{}{"http://127.0.0.1:9200/logc-test-"},
				"index":     "app-err",
				"mapping":   mappingIfc,
				"timeField": "at",
			},
			expect: testES3,
		},
	}
}
