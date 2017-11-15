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

var conf = config.Get()
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
	p.deleteObsolete()
	return true
}

// http://log-es.wumart.com/_cat/indices/logc-dev-*?h=index&s=index:desc
func (p *Pusher) deleteObsolete() {
	esAddr := elasticSearch.BaseAddrs[0]
	u, err := url.Parse(esAddr)
	if err != nil {
		p.logger.Errorf("parse es addr %s error: %+v\n", esAddr, err)
		return
	}
	for _, index := range p.getObsoleteIndices(u) {
		p.deleteIndex(u.Host, u.Scheme, index)
	}
}

func (p *Pusher) getObsoleteIndices(u *url.URL) []string {
	// u.path: /logc-dev-
	u.Path = fmt.Sprintf("/_cat/indices%s%s*", u.Path, p.file.Index)
	u.RawQuery = `h=index&s=index:desc`
	uri := u.String()
	res, err := httputil.Get(uri, nil, nil)
	if err != nil {
		p.logger.Errorf("get es history index %s error: %+v\n", uri, err)
		return nil
	}
	b, err := res.GetBody()
	if err != nil {
		p.logger.Errorf("get es history index error: ", err)
		return nil
	}
	esIndices := strings.Split(string(b), "\n")
	if len(esIndices) == 0 {
		p.logger.Error("no esIndices for ", uri)
		return nil
	}
	if len(esIndices) > p.file.IndexKeep {
		return esIndices[p.file.IndexKeep:]
	}
	return nil
}

func (p *Pusher) deleteIndex(host, scheme, index string) {
	u := url.URL{Host: host, Scheme: scheme, Path: index}
	uri := u.String()
	_, err := httputil.Delete(uri, nil, nil)
	if err != nil {
		p.logger.Errorf("delete es index %s error: %+v\n", uri, err)
	}
}
