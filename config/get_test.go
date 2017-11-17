package config

import (
	"os"
	"strings"
	"testing"

	"github.com/lovego/deep"
	"github.com/lovego/logc/collector/reader"
)

func TestGet(t *testing.T) {
	os.Args = []string{os.Args[0], `testdata/logc.yml`}
	got := Get()
	expect := getTestExpectConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Fatal("\n" + strings.Join(diff, "\n"))
	}
}

func getTestExpectConfig() Config {
	mapping := map[interface{}]interface{}{
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

	return Config{
		Name:    "test",
		Mailer:  "mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh",
		Keepers: []string{},
		Batch:   reader.Batch{Size: 102400, Wait: "3s"},
		Rotate: Rotate{
			Time: "33 8 1 * * *",
			Cmd:  []string{"logrotate", "logrotate.conf"},
		},
		Files: []File{
			{
				Path: "app.log",
				Outputs: []map[string]interface{}{{
					"@type":      "elastic-search",
					"addrs":      []interface{}{"http://log-es.wumart.com/logc-dev-"},
					"index":      "app-<2006.01.02>",
					"indexKeep":  3,
					"type":       "app-log",
					"mapping":    mapping,
					"timeField":  "at",
					"timeFormat": "2006-01-02T15:04:05Z0700",
				}},
			},
			{
				Path: "app.err",
				Outputs: []map[string]interface{}{{
					"@type":     "elastic-search",
					"addrs":     []interface{}{"http://log-es.wumart.com/logc-dev-"},
					"index":     "app-<2006.01.02>-err",
					"type":      "app-err",
					"mapping":   mapping,
					"timeField": "at",
				}},
			},
			{
				Path: "consume.log",
				Outputs: []map[string]interface{}{{
					"@type": "elastic-search",
					"addrs": []interface{}{"http://log-es.wumart.com/logc-dev-"},
					"index": "test-consume",
					"type":  "consume-log",
					"mapping": map[interface{}]interface{}{
						"at":   map[interface{}]interface{}{"type": "date"},
						"data": map[interface{}]interface{}{"type": "object"},
					},
					"timeField": "at",
				}},
			},
		},
	}
}
