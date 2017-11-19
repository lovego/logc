package pusher

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
)

var elasticSearch = elastic.New2(
	&httputil.Client{Client: http.DefaultClient}, conf.ElasticSearch...,
)

// TODO: more errors retry
func (p *Pusher) bulkCreate(index string, docs [][2]interface{}) [][2]interface{} {
	if errs := elasticSearch.BulkCreate(index+`/`+p.file.Type, docs); errs != nil {
		if err, ok := errs.(elastic.BulkError); ok {
			return err.FailedItems()
		}
		p.logger.Error("push err is not elastic.BulkError type, but %T", errs)
	}
	return nil
}

func (p *Pusher) ensureIndex(index string) bool {
	if err := elasticSearch.Ensure(index, nil); err != nil {
		p.logger.Errorf("ensure index %s error: %+v\n", index, err)
		return false
	}
	if err := elasticSearch.Put(index+`/_mapping/`+p.file.Type, map[string]interface{}{
		`properties`: p.file.Mapping,
	}, nil); err != nil {
		p.logger.Errorf("put mapping %s/%s error: %+v\n", index, p.file.Type, err)
		return false
	}
	p.deleteObsoleteIndices(index)
	return true
}

// http://log-es.wumart.com/_cat/indices/logc-dev-*?h=index&s=index:desc
func (p *Pusher) deleteObsoleteIndices(currentIndex string) {
	if p.file.IndexKeep <= 0 {
		return
	}
	for _, index := range p.getObsoleteIndices() {
		if index != currentIndex {
			if err := elasticSearch.Delete(index, nil); err != nil {
				p.logger.Errorf("delete index %s error: %s", index, err)
			}
		}
	}
}

func (p *Pusher) getObsoleteIndices() []string {
	indices := p.getIndices()
	if len(indices) <= p.file.IndexKeep {
		return nil
	}
	return indices[p.file.IndexKeep:]
}

func (p *Pusher) getIndices() (indices []string) {
	uri, err := url.Parse(conf.ElasticSearch[0])
	if err != nil {
		p.logger.Errorf("parse es addr %s error: %v", conf.ElasticSearch[0], err)
		return
	}

	indicesData := p.catIndices(*uri)
	if len(indicesData) == 0 {
		return
	}
	for _, data := range indicesData {
		index := strings.TrimPrefix(data.Index, uri.Path)
		if p.file.TimeSeriesIndex.Match(index) {
			indices = append(indices, index)
		}
	}
	if len(indices) == 0 {
		p.logger.Errorf("no indices matches: %s", p.file.Index)
	}
	return
}

type indexData struct {
	Index string `json:"index"`
}

func (p *Pusher) catIndices(uri url.URL) (result []indexData) {
	// uri.Path: /logc-dev-
	pattern := uri.Path + p.file.TimeSeriesIndex.Pattern()
	uri.Path = fmt.Sprintf("/_cat/indices%s", pattern)
	uri.RawQuery = `format=json&h=index&s=index:desc`
	uriStr := uri.String()

	if err := httputil.GetJson(uriStr, nil, nil, &result); err != nil {
		p.logger.Errorf("GET %s error: %+v\n", uriStr, err)
		return
	}
	if len(result) == 0 {
		p.logger.Errorf("no indices matches: %s", pattern)
	}
	return
}
