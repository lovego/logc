package time_series_index

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseTimeSeriesIndex(t *testing.T) {
	type testCaseT struct {
		input     string
		expect    *timeSeriesIndex
		expectErr error
	}
	testCases := []testCaseT{
		{``, nil, nil},
		{`test-index`, nil, nil},
		{`test-<2006.01`, nil, errors.New("invalid index: test-<2006.01")},
		{`test-<2006.01>-<index>`, nil, errors.New("invalid index: test-<2006.01>-<index>")},
		{`<2006.01>`, &timeSeriesIndex{timeLayout: `2006.01`}, nil},
		{`test-<2006.01>`, &timeSeriesIndex{prefix: `test-`, timeLayout: `2006.01`}, nil},
		{`test-<2006.01>-log`, &timeSeriesIndex{prefix: `test-`, timeLayout: `2006.01`, suffix: `-log`}, nil},
	}
	for _, tc := range testCases {
		got, err := parseTimeSeriesIndex(tc.input)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Fatalf("input: %s, expect: %v, got: %v\n", tc.input, tc.expect, got)
		}
		if !reflect.DeepEqual(err, tc.expectErr) {
			t.Fatalf("input: %s, expectErr: %v, got: %v\n", tc.input, tc.expectErr, err)
		}
	}
}

func TestTimeSeriesIndexMatch(t *testing.T) {
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
	tsi := timeSeriesIndex{prefix: `test-`, timeLayout: `2006.01.02`, suffix: `-log`}
	for _, tc := range testCases {
		got := tsi.Match(tc.input)
		if got != tc.expect {
			t.Fatalf("input: %s, expect: %v, got: %v\n", tc.input, tc.expect, got)
		}
	}
}
