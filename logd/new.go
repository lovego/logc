package logd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Logd struct {
	addr, mergeJson string
}

func New(addr, mergeJson string) (*Logd, error) {
	var err error
	if addr, err = parseAddress(addr); err != nil {
		return nil, err
	}
	if mergeJson != `` {
		if mergeJson, err = parseMergeJson(mergeJson); err != nil {
			return nil, err
		}
	}
	return &Logd{addr, mergeJson}, nil
}

func parseAddress(addr string) (string, error) {
	if addr == `` {
		return ``, errors.New(`logd address required.`)
	}
	if !strings.HasPrefix(addr, `http://`) && !strings.HasPrefix(addr, `https://`) {
		addr = `http://` + addr
	}
	return addr, nil
}

func parseMergeJson(merge string) (string, error) {
	mergeData := map[string]interface{}{}
	var err error
	if err = json.Unmarshal([]byte(merge), &mergeData); err == nil {
		var mergeJson []byte
		if mergeJson, err = json.Marshal(mergeData); err == nil {
			return string(mergeJson), nil
		}
	}
	return ``, fmt.Errorf("invalid merge json(%s): %s", err, merge)
}
