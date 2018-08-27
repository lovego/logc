package config

import (
	"os"
	"strings"
	"testing"

	"github.com/lovego/deep"
	"github.com/lovego/logc/collector/reader"
)

func TestGet(t *testing.T) {
	os.Args = []string{os.Args[0], `../testdata/logc.yml`}
	got := Get()
	expect := getTestExpectConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Fatal("\n" + strings.Join(diff, "\n"))
	}
}

func getTestExpectConfig() Config {
	mapping := map[interface{}]interface{}{
		"properties": map[interface{}]interface{}{
			"at":       map[interface{}]interface{}{"type": "date"},
			"duration": map[interface{}]interface{}{"type": "float"},
			"host":     map[interface{}]interface{}{"type": "keyword"},
			"method":   map[interface{}]interface{}{"type": "keyword"},
			"path": map[interface{}]interface{}{
				"type": "text",
				"fields": map[interface{}]interface{}{
					"raw": map[interface{}]interface{}{"type": "keyword"},
				},
			},
			"query":         map[interface{}]interface{}{"type": "object"},
			"rawQuery":      map[interface{}]interface{}{"type": "keyword"},
			"status":        map[interface{}]interface{}{"type": "keyword"},
			"req_body":      map[interface{}]interface{}{"type": "text"},
			"res_body":      map[interface{}]interface{}{"type": "text"},
			"req_body_size": map[interface{}]interface{}{"type": "integer"},
			"res_body_size": map[interface{}]interface{}{"type": "integer"},
			"ip":            map[interface{}]interface{}{"type": "ip"},
			"refer":         map[interface{}]interface{}{"type": "text"},
			"agent":         map[interface{}]interface{}{"type": "text"},
			"proto":         map[interface{}]interface{}{"type": "keyword"},
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

	return Config{
		Name:    "test_dev",
		Mailer:  "mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh",
		Keepers: []string{},
		Batch:   reader.Batch{Size: 102400, Wait: "3s"},
		Rotate: Rotate{
			Time: "33 8 1 * * *",
			Cmd:  []string{"logrotate", "logrotate.conf"},
		},
		Files: map[string]map[string]map[string]interface{}{
			"app.log": {
				"es": {
					"@type":         "elasticsearch",
					"addrs":         []interface{}{"http://127.0.0.1:9200/logc-dev-"},
					"index":         "app-log-<2006.01.02>",
					"indexKeep":     100,
					"mapping":       mapping,
					"timeField":     "at",
					"timeFormat":    "2006-01-02T15:04:05Z0700",
					"addTypeSuffix": true,
				},
			},
			"app.err": {
				"es": {
					"@type":         "elasticsearch",
					"addrs":         []interface{}{"http://127.0.0.1:9200/logc-dev-"},
					"index":         "app-err-<2006.01.02>",
					"mapping":       mapping,
					"timeField":     "at",
					"addTypeSuffix": true,
				},
			},
			"consume.log": {
				"es": {
					"@type": "elasticsearch",
					"addrs": []interface{}{"http://127.0.0.1:9200/logc-dev-"},
					"index": "consume-log",
					"mapping": map[interface{}]interface{}{
						"properties": map[interface{}]interface{}{
							"at":   map[interface{}]interface{}{"type": "date"},
							"data": map[interface{}]interface{}{"type": "object"},
						},
					},
					"timeField":     "at",
					"addTypeSuffix": true,
				},
			},
		},
	}
}
