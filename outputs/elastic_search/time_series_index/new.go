package time_series_index

import (
	"errors"
	"regexp"
	"strings"
	"time"

	loggerpkg "github.com/lovego/logger"
)

type TimeSeriesIndex struct {
	prefix     string
	timeLayout string
	suffix     string

	timeField  string
	timeFormat string
	keep       int
	logger     *loggerpkg.Logger
}

var timeSeriesIndexRegexp = regexp.MustCompile(`^([^<>]*)<([^<>]+)>([^<>]*)$`)

func New(index, timeField, timeFormat string, keep int, logger *loggerpkg.Logger) (
	*TimeSeriesIndex, error,
) {
	if strings.IndexByte(index, '<') < 0 && strings.IndexByte(index, '>') < 0 {
		return nil, nil
	}
	m := timeSeriesIndexRegexp.FindStringSubmatch(index)
	if len(m) != 4 {
		return nil, errors.New("invalid index: " + index)
	}

	if timeField == `` {
		return nil, errors.New(`empty timeField.`)
	}

	if timeFormat == `` {
		timeFormat = time.RFC3339
	}

	return &TimeSeriesIndex{
		prefix:     m[1],
		timeLayout: m[2],
		suffix:     m[3],
		timeField:  timeField,
		timeFormat: timeFormat,
		keep:       keep,
		logger:     logger,
	}, nil
}
