package time_series_index

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lovego/xiaomei/utils/elastic"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

var testES = elastic.New(`http://log-es.com/logc-test-`)
var testIndices = []string{
	`not-log-index`, `log-2017xxx`, `log-2017.13`,
	`log-2017.12`, `log-2017.11`, `log-2017.10`, `log-2017.09`, `log-2017.08`, `log-2017.07`,
	`log-2017.06`, `log-2017.05`, `log-2017.04`, `log-2017.03`, `log-2017.02`, `log-2017.01`,
}

var testTimeSeriesIndex = TimeSeriesIndex{
	prefix: `log-`, timeLayout: `2006.01`, logger: loggerpkg.New(``, os.Stderr, nil), keep: 6,
}

func TestPrune(t *testing.T) {
	ensureTestIndices()
	testTimeSeriesIndex.Prune(testES)
	got := getExistingTestIndices()
	expect := testIndices[:9]
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("expect: %v, got: %v\n", expect, got)
	}
}

func TestGetIndices(t *testing.T) {
	ensureTestIndices()
	got := testTimeSeriesIndex.getIndices(testES)
	expect := testIndices[3:]
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("expect: %v, got: %v\n", expect, got)
	}
}

func TestCatIndices(t *testing.T) {
	ensureTestIndices()
	got := testTimeSeriesIndex.catIndices(testES)
	expect := testIndices[1:]
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("expect: %v, got: %v\n", expect, got)
	}
}

func TestMatch(t *testing.T) {
	type testCaseT struct {
		input  string
		expect bool
	}
	testCases := []testCaseT{
		{``, false},
		{`test-log`, false},
		{`test-2017.11.16-log`, true},
		{`test-2017.11.00-log`, false},
	}
	tsi := TimeSeriesIndex{prefix: `test-`, timeLayout: `2006.01.02`, suffix: `-log`}
	for _, tc := range testCases {
		got := tsi.match(tc.input)
		if got != tc.expect {
			t.Fatalf("input: %s, expect: %v, got: %v\n", tc.input, tc.expect, got)
		}
	}
}

func ensureTestIndices() {
	existing := getExistingTestIndices()

	m := map[string]bool{}
	for _, index := range existing {
		m[index] = true
	}

	for _, index := range testIndices {
		if !m[index] {
			if err := testES.Ensure(index, nil); err != nil {
				log.Panic(err)
			}
		}
		delete(m, index)
	}
	for index := range m {
		if err := testES.Delete(index, nil); err != nil {
			log.Panic(err)
		}
	}
}

func getExistingTestIndices() (result []string) {
	slice := []struct {
		Index string `json:"index"`
	}{}
	urlStr := "/_cat/indices/logc-test-*?format=json&h=index&s=index:desc"
	if err := testES.RootGet(urlStr, nil, &slice); err != nil {
		log.Panicf("GET %s error: %+v\n", urlStr, err)
	}
	for _, one := range slice {
		result = append(result, strings.TrimPrefix(one.Index, `logc-test-`))
	}
	return
}
