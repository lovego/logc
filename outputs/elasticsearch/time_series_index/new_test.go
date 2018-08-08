package time_series_index

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/lovego/deep"
	loggerpkg "github.com/lovego/logger"
)

var testLogger = loggerpkg.New(os.Stderr).SetAlarm(nil)

func init() {
	deep.CompareUnexportedFields = true
}

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
		got, err := New(``, tc.input, `at`, ``, 0, testLogger)
		expect := tc.expect
		if expect != nil {
			expect.timeField = `at`
			expect.timeFormat = time.RFC3339
			expect.logger = testLogger
		}
		if diff := deep.Equal(got, expect); diff != nil {
			t.Fatalf("\ninput: %s\n diff: %+v\n", tc.input, diff)
		}
		if diff := deep.Equal(err, tc.expectErr); diff != nil {
			t.Fatalf("\ninput: %s\n diff: %+v\n", tc.input, diff)
		}
	}
}
