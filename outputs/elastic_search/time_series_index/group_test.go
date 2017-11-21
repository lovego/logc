package time_series_index

import (
	"os"
	"reflect"
	"testing"
	"time"

	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

func TestGroup(t *testing.T) {
	type testCaseT struct {
		input  []map[string]interface{}
		expect []Rows
	}
	testCases := []testCaseT{
		{
			input: []map[string]interface{}{
				{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v1`},
				{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v2`},
				{`at`: `2017-11-01T00:00:00+08:00`, `k`: `v3`},
			},
			expect: []Rows{
				{
					Index: `test-2017.10.31-log`, Rows: []map[string]interface{}{
						{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v1`},
						{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v2`},
					},
				},
				{
					Index: `test-2017.11.01-log`, Rows: []map[string]interface{}{
						{`at`: `2017-11-01T00:00:00+08:00`, `k`: `v3`},
					},
				},
			},
		},
		{
			input: []map[string]interface{}{
				{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v1`},
				{`at`: `2017-10-31T23:59:59+08:00`, `k`: `v2`},
				{`k`: `v3`},
			},
		},
	}
	tsi := TimeSeriesIndex{
		prefix: `test-`, timeLayout: `2006.01.02`, suffix: `-log`,
		timeField: `at`, timeFormat: time.RFC3339, logger: loggerpkg.New(``, os.Stderr, nil),
	}
	for _, tc := range testCases {
		got := tsi.Group(tc.input)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("input: %s, expect: %v, got: %v\n", tc.input, tc.expect, got)
		}
	}
}
