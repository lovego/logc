package elasticsearch

import (
	"strings"

	"reflect"

	"log"

	"math"

	"github.com/nu7hatch/gouuid"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (es *ElasticSearch) Write(rows []map[string]interface{}) bool {
	if len(rows) == 0 {
		return true
	}
	if es.timeSeriesIndex == nil {
		es.writeToIndex(es.index, rows)
		return true
	}

	return es.writeToTimeSeriesIndex(rows)
}

func (es *ElasticSearch) writeToTimeSeriesIndex(rows []map[string]interface{}) bool {
	indicesRows := es.timeSeriesIndex.Group(rows)
	if len(indicesRows) <= 0 { // Group encountered error.
		return false
	}
	prune := false
	for _, one := range indicesRows {
		if one.Index != es.currentIndex {
			if !es.ensureIndex(one.Index, es.logger) {
				return false
			}
			es.currentIndex = one.Index
			prune = true
		}
		es.writeToIndex(one.Index, one.Rows)
	}
	if prune {
		es.timeSeriesIndex.Prune(es.client)
	}
	return true
}

func (es *ElasticSearch) writeToIndex(index string, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}
	if es.addTypeSuffix {
		rows = es.structure(convert2InterfaceSlice(rows), false)
	}
	docs := es.addDocId(rows)

	var t Timer
	for {
		if docs = es.bulkCreate(index, docs); len(docs) == 0 {
			break
		}
		t.Sleep()
	}
}

func (es *ElasticSearch) addDocId(rows []map[string]interface{}) [][2]interface{} {
	docs := [][2]interface{}{}
	for _, doc := range rows {
		convertKeyWithDot(doc)
		if uid, err := uuid.NewV4(); err != nil {
			es.logger.Errorf("generate uuid error: %v", err)
			docs = append(docs, [2]interface{}{nil, doc})
		} else {
			docs = append(docs, [2]interface{}{strings.Replace(uid.String(), `-`, ``, -1), doc})
		}
	}
	return docs
}

// convert dot(.) in key to underline(_)
func convertKeyWithDot(doc map[string]interface{}) {
	for key, value := range doc {
		if strings.ContainsRune(key, '.') {
			newKey := strings.Replace(key, `.`, `_`, -1)
			doc[newKey] = value
			delete(doc, key)
		}
		if v, ok := value.(map[string]interface{}); ok {
			convertKeyWithDot(v)
		}
		if vs, ok := value.([]map[string]interface{}); ok {
			for _, v := range vs {
				convertKeyWithDot(v)
			}
		}
	}
}

func convert2InterfaceSlice(s []map[string]interface{}) []interface{} {
	m := make([]interface{}, 0)

	for _, v := range s {
		m = append(m, interface{}(v))
	}

	return m
}

func assertInterfaceSliceSuffix(slice []interface{}) string {
	if len(slice) == 0 {
		return ""
	}
	types := map[string]string{
		"string": "_s",

		"int8":  "_i",
		"uint8": "_i",

		"int16":  "_i",
		"uint16": "_i",

		"int32":  "_i",
		"uint32": "_i",

		"int":  "_i",
		"uint": "_i",

		"int64":  "_i",
		"uint64": "_i",

		"float32": "_f",
		"float64": "_f",
	}
	var typo string
	typo = reflect.TypeOf(slice[0]).String()
	for _, v := range slice {
		if typo != reflect.TypeOf(v).String() {
			return ""
		}
	}
	return types[typo]
}

func suuffixForESNumberType(f float64) string {
	if math.Trunc(f) == f {
		return "_i"
	} else {
		return "_f"
	}
}

func (es *ElasticSearch) structure(rows []interface{}, recursive bool) []map[string]interface{} {
	rewrited := make([]map[string]interface{}, 0)
	for _, row := range rows {
		newRow := make(map[string]interface{}, 0)
		var newKey string
		var newValue interface{}
		if assertedRow, ok := row.(map[string]interface{}); ok {
			for oldKey, value := range assertedRow {
				if es.has(oldKey) && !recursive {
					newKey = oldKey
					newValue = value
				} else {
					switch assertedValue := value.(type) {
					case string:
						newKey = oldKey + "_s"
						newValue = value
					case int8, uint8, int16, uint16, int, uint, int32, uint32, int64, uint64:
						newKey = oldKey + "_i"
						newValue = value
					case float32, float64:
						newKey = oldKey + suuffixForESNumberType(assertedValue.(float64))
						newValue = value
					case bool:
						newKey = oldKey + "_b"
						newValue = value
					case []map[string]interface{}:
						newKey = oldKey + "_o"
						newValue = es.structure(convert2InterfaceSlice(assertedValue), true)
					case map[string]interface{}:
						newKey = oldKey + "_o"
						newValue = es.structure(convert2InterfaceSlice([]map[string]interface{}{assertedValue}), true)
						if nv, ok := newValue.([]map[string]interface{}); ok && len(nv) > 0 {
							newValue = nv[0]
						}
					case []string:
						newKey = oldKey + "_s"
						newValue = value
					case []int, []uint, []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64:
						newKey = oldKey + "_i"
						newValue = value
					case []interface{}:
						suffix := assertInterfaceSliceSuffix(assertedValue)
						if suffix != "" {
							if suffix == "_f" {
								newKey = oldKey + suuffixForESNumberType(assertedValue[0].(float64))
							} else {
								newKey = oldKey + suffix
							}
							newValue = assertedValue
						} else {
							newKey = oldKey + "_o"
							newValue = es.structure(assertedValue, true)
							if len(newValue.([]map[string]interface{})) == 0 {
								newKey = oldKey + "_a"
								newValue = assertedValue
							}
						}
					default:
						continue
					}
				}
				if newKey != "" {
					newRow[newKey] = newValue
				}
			}
		}
		if len(newRow) != 0 {
			rewrited = append(rewrited, newRow)
		}
	}
	return rewrited
}

func (es *ElasticSearch) has(property string) bool {
	if properties, ok := es.mapping["properties"]; ok {
		if values, ok := properties.(map[string]interface{}); ok {
			_, ok = values[property]
			return ok
		}
	}
	return false
}
