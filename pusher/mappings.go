package pusher

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
)

var dataEs *elastic.ES

func (p *Pusher) ensureIndex(esIndex string) {
	if err := dataEs.Ensure(esIndex, nil); err != nil {
		p.logger.Fatalf("create files error: %+v\n", err)
	}
	if err := dataEs.Put(esIndex+`/_mapping/`+p.Type, map[string]interface{}{
		`properties`: p.Mapping,
	}, nil); err != nil {
		p.logger.Fatalf("create files error: %+v\n", err)
	}
	p.delHistory()
}

// http://log-es.wumart.com/_cat/indices/logc-dev-*?h=index&s=index:desc
func (p *Pusher) delHistory() {
	esAddr := dataEs.BaseAddrs[0]
	u, err := url.Parse(esAddr)
	if err != nil {
		p.logger.Errorf("parse es addr %s error: %+v\n", esAddr, err)
		return
	}
	for _, esIndex := range p.indicesToDel(u) {
		p.deleteIndex(u.Host, u.Scheme, esIndex)
	}
}

func (p *Pusher) indicesToDel(u *url.URL) []string {
	// u.path: /logc-dev-
	u.Path = fmt.Sprintf("/_cat/indices%s%s*", u.Path, p.Index)
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
	if len(esIndices) > p.Keep {
		return esIndices[p.Keep:]
	}
	return nil
}

func (p *Pusher) deleteIndex(host, scheme, esIndex string) {
	u := url.URL{Host: host, Scheme: scheme, Path: esIndex}
	uri := u.String()
	_, err := httputil.Delete(uri, nil, nil)
	if err != nil {
		p.logger.Errorf("delete es index %s error: %+v\n", uri, err)
	}
}
