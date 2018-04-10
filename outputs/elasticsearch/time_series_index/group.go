package time_series_index

import (
	"time"
)

type Rows struct {
	Index string
	Rows  []map[string]interface{}
}

func (tsi *TimeSeriesIndex) Group(rows []map[string]interface{}) (result []Rows) {
	indices := []string{}
	m := make(map[string][]map[string]interface{})
	for _, row := range rows {
		if index := tsi.Of(row); index != `` {
			if m[index] == nil {
				m[index] = []map[string]interface{}{row}
				indices = append(indices, index)
			} else {
				m[index] = append(m[index], row)
			}
		} else {
			return nil
		}
	}
	for _, index := range indices {
		result = append(result, Rows{Index: index, Rows: m[index]})
	}
	return
}

func (tsi TimeSeriesIndex) Of(row map[string]interface{}) string {
	value, ok := row[tsi.timeField].(string)
	if !ok {
		tsi.logger.Errorf("%s: non string timeField %s: %v",
			tsi.collectorId, tsi.timeField, row[tsi.timeField],
		)
		return ``
	}
	at, err := time.Parse(tsi.timeFormat, value)
	if err != nil {
		tsi.logger.Errorf("%s: parse timeField %s with layout %s error: %v",
			tsi.collectorId, tsi.timeField, tsi.timeFormat, err,
		)
		return ``
	}
	return tsi.Get(at)
}

func (tsi TimeSeriesIndex) Get(t time.Time) string {
	return tsi.prefix + t.Format(tsi.timeLayout) + tsi.suffix
}
