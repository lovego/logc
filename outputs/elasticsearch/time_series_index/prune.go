package time_series_index

import (
	"net/url"
	"strings"
	"time"

	"github.com/lovego/elastic"
)

func (tsi TimeSeriesIndex) Prune(client *elastic.ES) {
	if tsi.keep <= 0 {
		return
	}
	indices := tsi.getIndices(client)
	if len(indices) <= tsi.keep {
		return
	}
	obsoletes := indices[tsi.keep:]
	for _, index := range obsoletes {
		if err := client.Delete(index, nil); err != nil && !elastic.IsNotFound(err) {
			tsi.logger.Errorf("%s: delete index %s error: %s", tsi.collectorId, index, err)
		} else {
			tsi.logger.Printf("pruned index: %s", index)
		}
	}
}

func (tsi TimeSeriesIndex) getIndices(client *elastic.ES) (indices []string) {
	indicesList := tsi.catIndices(client)
	if len(indicesList) == 0 {
		return
	}
	for _, index := range indicesList {
		if tsi.match(index) {
			indices = append(indices, index)
		}
	}
	if len(indices) == 0 {
		tsi.logger.Errorf("%s: no indices matches: %s<%s>%s",
			tsi.collectorId, tsi.prefix, tsi.timeLayout, tsi.suffix,
		)
	}
	return
}

func (tsi TimeSeriesIndex) catIndices(client *elastic.ES) (indices []string) {
	uri, err := url.Parse(client.BaseAddrs[0])
	if err != nil {
		tsi.logger.Errorf("%s: parse es addr %s error: %v", tsi.collectorId, client.BaseAddrs[0], err)
		return
	}
	pattern := uri.Path + tsi.prefix + `*` + tsi.suffix // uri.Path: /logc-dev-

	var result []struct {
		Index string `json:"index"`
	}
	query := "/_cat/indices" + pattern + "?format=json&h=index&s=index:desc"
	if err := client.RootGet(query, nil, &result); err != nil {
		tsi.logger.Errorf("%s: %s error: %+v\n", tsi.collectorId, query, err)
		return
	}
	if len(result) == 0 {
		tsi.logger.Errorf("%s: no indices matches: %s", tsi.collectorId, pattern)
	}
	prefix := strings.TrimPrefix(uri.Path, `/`)
	for _, data := range result {
		indices = append(indices, strings.TrimPrefix(data.Index, prefix))
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
