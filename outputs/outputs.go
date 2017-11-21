package outputs

import (
	"github.com/lovego/logc/outputs/elastic_search"
	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

var logger *loggerpkg.Logger

func Setup(l *loggerpkg.Logger) {
	logger = l
}

type Output interface {
	Write(rows []map[string]interface{}) (ok bool)
}

// Different collector must use separate output. Because output has internal state.
// For example. elastic_search has currentIndex state, it is designed only one file in mind.
// So, when a collector is constucted, use a maker to make a new ouput.
func Maker(conf map[string]interface{}, file string) func(*loggerpkg.Logger) Output {
	typ := getType(conf, file)
	if typ == `` {
		return nil
	}
	return func(logger *loggerpkg.Logger) Output {
		return New(typ, conf, file, logger)
	}
}

func New(typ string, conf map[string]interface{}, file string, logger *loggerpkg.Logger) Output {
	switch typ {
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

func getType(conf map[string]interface{}, file string) string {
	typeV := conf[`@type`]
	if typeV == nil {
		logger.Fatalf("%s: no @type defined.", file)
		return ``
	}
	typ, ok := typeV.(string)
	if !ok {
		logger.Fatalf("%s: non string @type defined.", file)
		return ``
	}
	switch typ {
	case `elastic-search`:
		return typ
	default:
		logger.Fatalf("%s: unknown @type: %s .", file, typ)
		return ``
	}
}
