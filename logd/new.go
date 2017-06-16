package logd

import (
	"encoding/json"
	"net/url"

	"github.com/lovego/xiaomei/utils"
)

type Logd struct {
	addr, mergeJson string
}

func New(addr, mergeJson string) *Logd {
	if addr = parseAddress(addr); addr == `` {
		return nil
	}
	if mergeJson != `` {
		if mergeJson = parseMergeJson(mergeJson); mergeJson == `` {
			return nil
		}
	}
	return &Logd{addr, mergeJson}
}

func parseAddress(addr string) string {
	if addr == `` {
		utils.Log(`logd address required.`)
		return ``
	}
	u, err := url.Parse(addr)
	if err != nil {
		utils.Logf("invalid logd address(%v): %s", err, addr)
		return ``
	}
	if u.Host == `` {
		utils.Logf("invalid logd address(%s): %s", `empty host`, addr)
		return ``
	}
	if u.Scheme == `` {
		u.Scheme = `http`
	}
	return (&url.URL{Scheme: u.Scheme, Host: u.Host}).String()
}

func parseMergeJson(merge string) string {
	mergeData := map[string]interface{}{}
	var err error
	if err = json.Unmarshal([]byte(merge), &mergeData); err == nil {
		var mergeJson []byte
		if mergeJson, err = json.Marshal(mergeData); err == nil {
			return string(mergeJson)
		}
	}
	utils.Logf("invalid merge json(%s): %s", err, merge)
	return ``
}
