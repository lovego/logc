package logd

import (
	"encoding/json"
	"strings"

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
	if !strings.HasPrefix(addr, `http://`) && !strings.HasPrefix(addr, `https://`) {
		addr = `http://` + addr
	}
	return addr
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
