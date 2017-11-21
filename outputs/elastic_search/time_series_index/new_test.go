package time_series_index

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type testCaseT struct {
		input     string
		expect    *TimeSeriesIndex
		expectErr error
	}
	testCases := []testCaseT{
		{``, nil, nil},
		{`test-index`, nil, nil},
		{`test-<2006.01`, nil, errors.New("invalid index: test-<2006.01")},
		{`test-<2006.01>-<index>`, nil, errors.New("invalid index: test-<2006.01>-<index>")},
		{`<2006.01>`, &TimeSeriesIndex{timeLayout: `2006.01`}, nil},
		{`test-<2006.01>`, &TimeSeriesIndex{prefix: `test-`, timeLayout: `2006.01`}, nil},
		{`test-<2006.01>-log`, &TimeSeriesIndex{
			prefix: `test-`, timeLayout: `2006.01`, suffix: `-log`,
		}, nil},
	}
	for _, tc := range testCases {
		got, err := New(tc.input, `at`, ``, 0, nil)
		expect := tc.expect
		if expect != nil {
			expect.timeField = `at`
			expect.timeFormat = time.RFC3339
		}
		if !reflect.DeepEqual(got, expect) {
			t.Fatalf("input: %s, expect: %v, got: %v\n", tc.input, tc.expect, got)
		}
		if !reflect.DeepEqual(err, tc.expectErr) {
			t.Fatalf("input: %s, expectErr: %v, got: %v\n", tc.input, tc.expectErr, err)
		}
	}
}
