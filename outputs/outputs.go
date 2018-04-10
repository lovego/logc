package outputs

import (
	"github.com/lovego/logc/outputs/elasticsearch"
	loggerpkg "github.com/lovego/logger"
)

var theLogger *loggerpkg.Logger

func Setup(logger *loggerpkg.Logger) {
	theLogger = logger
	elasticsearch.Setup(logger)
}

type Output interface {
	Write(rows []map[string]interface{}) (ok bool)
}

// Different collector must use separate output. Because output has internal state.
// For example, elasticsearch has currentIndex state, it's designed having only one file in mind.
// So, when a collector is constructed, use a maker to make a new ouput.
func Maker(collectorId string, conf map[string]interface{}) func(*loggerpkg.Logger) Output {
	typ := getType(conf, collectorId)
	if typ == `` {
		return nil
	}
	return func(logger *loggerpkg.Logger) Output {
		return New(collectorId, typ, conf, logger)
	}
}

func New(collectorId string, typ string, conf map[string]interface{}, logger *loggerpkg.Logger) Output {
	switch typ {
	case `elasticsearch`:
		// the if is required. because nil pointer makes a non nil interface.
		if output := elasticsearch.New(collectorId, conf, logger); output != nil {
			return output
		} else {
			return nil
		}
	default:
		theLogger.Errorf("unknown output @type: %v", typ)
		return nil
	}
}

func getType(conf map[string]interface{}, collectorId string) string {
	typeV := conf[`@type`]
	if typeV == nil {
		theLogger.Fatalf("%s: no @type defined.", collectorId)
		return ``
	}
	typ, ok := typeV.(string)
	if !ok {
		theLogger.Fatalf("%s: non string @type defined.", collectorId)
		return ``
	}
	switch typ {
	case `elasticsearch`:
		return typ
	default:
		theLogger.Fatalf("%s: unknown @type: %s .", collectorId, typ)
		return ``
	}
}
