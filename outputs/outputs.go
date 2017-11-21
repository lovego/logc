package outputs

import (
	"github.com/lovego/logc/outputs/elastic_search"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

type Output interface {
	Write(rows []map[string]interface{}) (ok bool)
}

// Different collector must use separate output. Because output has internal state.
// For example. elastic_search has currentIndex state, it is designed only one file in mind.
// So, when a collector is constucted, use a maker to make a new ouput.
func Maker(conf map[string]interface{}, file string) func(*loggerpkg.Logger) Output {
	return func(logger *loggerpkg.Logger) Output {
		return New(conf, file, logger)
	}
}

func New(conf map[string]interface{}, file string, logger *loggerpkg.Logger) Output {
	switch typ := conf[`@type`].(string); typ {
	case `elastic-search`:
		// the if is required. because nil pointer makes a non nil interface.
		if output := elastic_search.New(conf, file, logger); output != nil {
			return output
		} else {
			return nil
		}
	default:
		logger.Errorf("unknown output @type: %v", typ)
		return nil
	}
}
