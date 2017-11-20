package time_series_index

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
)

func (tsi TimeSeriesIndex) Prune(client *elastic.ES) {
	if tsi.keep <= 0 {
		return
	}
	indices := tsi.getIndices(client.BaseAddrs[0])
	if len(indices) <= tsi.keep {
		return
	}
	obsoletes := indices[tsi.keep:]
	for _, index := range obsoletes {
		if err := client.Delete(index, nil); err != nil {
			tsi.logger.Errorf("delete index %s error: %s", index, err)
		}
	}
}

func (tsi TimeSeriesIndex) getIndices(baseAddr string) (indices []string) {
	uri, err := url.Parse(baseAddr)
	if err != nil {
		tsi.logger.Errorf("parse es addr %s error: %v", baseAddr, err)
		return
	}

	indicesData := tsi.catIndices(*uri)
	if len(indicesData) == 0 {
		return
	}
	for _, data := range indicesData {
		index := strings.TrimPrefix(data.Index, uri.Path)
		if tsi.match(index) {
			indices = append(indices, index)
		}
	}
	if len(indices) == 0 {
		tsi.logger.Errorf("no indices matches: %s<%s>%s", tsi.prefix, tsi.timeLayout, tsi.suffix)
	}
	return
}

func (tsi TimeSeriesIndex) match(index string) bool {
	if len(index) <= len(tsi.prefix)+len(tsi.suffix) ||
		!strings.HasPrefix(index, tsi.prefix) ||
		!strings.HasSuffix(index, tsi.suffix) {
		return false
	}
	timeStr := index[len(tsi.prefix) : len(index)-len(tsi.suffix)]
	_, err := time.Parse(tsi.timeLayout, timeStr)
	return err == nil
}

type indexData struct {
	Index string `json:"index"`
}

func (tsi TimeSeriesIndex) catIndices(uri url.URL) (result []indexData) {
	// uri.Path: /logc-dev-
	pattern := uri.Path + tsi.prefix + `*` + tsi.suffix

	// http://log-es.wumart.com/_cat/indices/logc-dev-*?h=index&s=index:desc
	uri.Path = fmt.Sprintf("/_cat/indices%s", pattern)
	uri.RawQuery = `format=json&h=index&s=index:desc`
	uriStr := uri.String()

	if err := httputil.GetJson(uriStr, nil, nil, &result); err != nil {
		tsi.logger.Errorf("GET %s error: %+v\n", uriStr, err)
		return
	}
	if len(result) == 0 {
		tsi.logger.Errorf("no indices matches: %s", pattern)
	}
	return
}
