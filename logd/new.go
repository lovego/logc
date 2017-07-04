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

func New(addr string, mergeData map[string]interface{}) (*Logd, error) {
	var err error
	if addr, err = parseAddress(addr); err != nil {
		return nil, err
	}
	var mergeJson []byte
	if len(mergeData) > 0 {
		if mergeJson, err = json.Marshal(mergeData); err != nil {
			return nil, fmt.Errorf("marshal merge data: %v", err)
		}
	}
	return &Logd{addr, string(mergeJson)}, nil
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
