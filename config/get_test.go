package config

import (
	//	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lovego/deep"
)

func TestGet(t *testing.T) {
	os.Args = []string{os.Args[0], `testdata/logc.yml`}
	got := Get()
	expect := getTestExpectConfig()

	if diff := deep.Equal(got, expect); diff != nil {
		t.Fatal(strings.Join(diff, "\n"))
	}
}

func getTestExpectConfig() Config {
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

	return Config{
		Name:              "test",
		ElasticSearch:     []string{"http://log-es.wumart.com/logc-dev-"},
		BatchSize:         102400,
		BatchWait:         "3s",
		BatchWaitDuration: 3 * time.Second,
		RotateTime:        "33 8 1 * * *",
		RotateCmd:         []string{"logrotate", "logrotate.conf"},
		Mailer:            "mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh",
		Keepers:           []string{},
		Files: []File{
			{
				Path:            "app.log",
				Index:           "app-<2006.01.02>",
				TimeSeriesIndex: &timeSeriesIndex{prefix: `app-`, timeLayout: `2006.01.02`},
				IndexKeep:       3,
				Type:            "app-log",
				Mapping:         mapping,
				TimeField:       "at",
				TimeFormat:      "2006-01-02T15:04:05Z0700",
			},
			{
				Path:            "app.err",
				Index:           "app-<2006.01.02>-err",
				TimeSeriesIndex: &timeSeriesIndex{prefix: `app-`, timeLayout: `2006.01.02`, suffix: `-err`},
				Type:            "app-err",
				Mapping:         mapping,
				TimeField:       "at",
				TimeFormat:      time.RFC3339,
			},
			{
				Path:  "consume.log",
				Index: "test-consume",
				Type:  "consume-log",
				Mapping: map[string]map[string]interface{}{
					"at":   {"type": "date"},
					"data": {"type": "object"},
				},
				TimeField:  "at",
				TimeFormat: time.RFC3339,
			},
		},
	}
}
