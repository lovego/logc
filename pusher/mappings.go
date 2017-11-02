package pusher

import (
	"net/http"

	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/elastic"
	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/lovego/xiaomei/utils/logger"
)

type File struct {
	Type    string                            `yaml:"name"`
	Mapping map[string]map[string]interface{} `yaml:"mapping"`
}

var dataEs *elastic.ES

func CreateMappings(esAddrs []string, filesAry []*config.File, log *logger.Logger) {
	if dataEs == nil {
		dataEs = elastic.New2(&httputil.Client{Client: http.DefaultClient}, esAddrs...)
	}
	mappings := make(map[string][]File)
	for _, file := range filesAry {
		index := file.Index
		if mappings[index] == nil {
			mappings[index] = []File{}
		}
		mappings[index] = append(mappings[index], File{file.Type, file.Mapping})
	}
	for index, files := range mappings {
		for _, file := range files {
			if err := dataEs.Ensure(index, nil); err != nil {
				log.Fatalf("create files error: %+v\n", err)
			}
			if err := dataEs.Put(index+`/_mapping/`+file.Type, map[string]interface{}{
				`properties`: file.Mapping,
			}, nil); err != nil {
				log.Fatalf("create files error: %+v\n", err)
			}
		}
	}
}
