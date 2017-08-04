package pusher

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/httputil"
)

func CreateMappings(logdAddr string, filesAry []*config.File) {
	mappings := make(map[string][]map[string]interface{})
	for _, file := range filesAry {
		org := file.Org
		if mappings[org] == nil {
			mappings[org] = []map[string]interface{}{}
		}
		mappings[org] = append(mappings[org], map[string]interface{}{
			`name`: file.Name, `mapping`: file.Mapping,
		})
	}
	for org, filesMapping := range mappings {
		createMappings(logdAddr, org, filesMapping)
	}
}

func createMappings(logdAddr, org string, files []map[string]interface{}) {
	filesJson, err := json.Marshal(files)
	if err != nil {
		log.Fatal("marshal mappings error: ", err)
	}
	query := url.Values{}
	query.Set(`org`, org)
	createUrl := logdAddr + `/org-files?` + query.Encode()
	resp := struct {
		Code, Message string
	}{}
	if err := httputil.PostJson(createUrl, nil, filesJson, &resp); err != nil {
		log.Fatal("create files error: %+v\n", err)
	}
	if resp.Code != `ok` {
		log.Fatal("create files failed: %+v\n", resp)
	}
}