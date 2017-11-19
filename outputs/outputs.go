package outputs

import (
	"github.com/lovego/logc/outputs/elastic_search"
	"github.com/lovego/xiaomei/utils/logger"
)

type Output interface {
	Write(rows []map[string]interface{}, logger *loggerpkg.Logger) bool
}

func Check(conf map[string]interface{}) {
}

func New(conf map[string]interface{}, file string) Output {
	switch typ := conf[`@type`].(string); typ {
	case `elastic-search`:
		if o := elastic_search.New(conf, file); o != nil {
			return o
		} else {
			return nil
		}
	default:
		logger.Errorf("unknown output @type: %v", typ)
		return nil
	}
}
