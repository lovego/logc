package elastic_search

import (
	"testing"
)

func TestWrite1(t *testing.T) {
	clearTestIndex(`app-2017.07.03`, t)
	clearTestIndex(`app-2017.08.04`, t)

	rows := []map[string]interface{}{
		{"at": "2017-07-03T08:49:16+0800", "duration": "0.000077", "host": "example.dev", "ip": "192.168.56.1", "method": "GET", "path": "/result-ok", "proto": "HTTP/1.0", "req_body": 0, "res_body": 60, "status": 200},
		{"at": "2017-07-03T20:49:16+0800", "duration": "0.000077", "host": "example.dev", "ip": "192.168.56.1", "method": "GET", "path": "/result-ok", "proto": "HTTP/1.0", "req_body": 0, "res_body": 60, "status": 200},
		{"at": "2017-07-03T20:49:16+0800", "duration": "0.000077", "host": "example.dev", "ip": "192.168.56.1", "method": "GET", "path": "/result-ok", "proto": "HTTP/1.0", "req_body": 0, "res_body": 60, "status": 200},
		{"at": "2017-08-04T08:49:16+0800", "duration": "0.000077", "host": "example.dev", "ip": "192.168.56.1", "method": "GET", "path": "/result-ok", "proto": "HTTP/1.0", "req_body": 0, "res_body": 60, "status": 200},
	}
	if !testES1.Write(rows) {
		t.Fatalf(`Write rows failed.`)
	}
	if expect, got := 3, getTestIndexDocCount(`app-2017.07.03`, t); got != expect {
		t.Errorf("app-2017.07.03: expect %d docs, got %d.\n", expect, got)
	}
	if expect, got := 1, getTestIndexDocCount(`app-2017.08.04`, t); got != expect {
		t.Errorf("app-2017.08.04: expect %d docs, got %d.\n", expect, got)
	}
}

func getTestIndexDocCount(index string, t *testing.T) int {
	if err := testEsClient.Get(index+`/_refresh`, nil, nil); err != nil {
		t.Fatal(err)
	}
	result := struct {
		Count int `json:"count"`
	}{}
	if err := testEsClient.Get(index+`/_count`, nil, &result); err != nil {
		t.Fatal(err)
	}
	return result.Count
}

func clearTestIndex(index string, t *testing.T) {
	if ok, err := testEsClient.Exist(index); err != nil {
		t.Fatal(err)
	} else if ok {
		if err := testEsClient.Delete(index, nil); err != nil {
			t.Fatal(err)
		}
	}
}
