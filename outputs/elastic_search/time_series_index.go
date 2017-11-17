package config

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type timeSeriesIndex struct {
	prefix     string
	timeLayout string
	suffix     string
}

var timeSeriesIndexRegexp = regexp.MustCompile(`^([^<>]*)<([^<>]+)>([^<>]*)$`)

func parseTimeSeriesIndex(index string) (*timeSeriesIndex, error) {
	if strings.IndexByte(index, '<') < 0 && strings.IndexByte(index, '>') < 0 {
		return nil, nil
	}
	m := timeSeriesIndexRegexp.FindStringSubmatch(index)
	if len(m) != 4 {
		return nil, errors.New("invalid index: " + index)
	}
	return &timeSeriesIndex{
		prefix:     m[1],
		timeLayout: m[2],
		suffix:     m[3],
	}, nil
}

func (i timeSeriesIndex) Get(t time.Time) string {
	return i.prefix + t.Format(i.timeLayout) + i.suffix
}

func (i timeSeriesIndex) Pattern() string {
	return i.prefix + `*` + i.suffix
}

func (i timeSeriesIndex) Match(index string) bool {
	if len(index) <= len(i.prefix)+len(i.suffix) ||
		!strings.HasPrefix(index, i.prefix) ||
		!strings.HasSuffix(index, i.suffix) {
		return false
	}
	timeStr := index[len(i.prefix) : len(index)-len(i.suffix)]
	_, err := time.Parse(i.timeLayout, timeStr)
	return err == nil
}
