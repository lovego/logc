package elasticsearch

import (
	"encoding/json"
	"strings"
)

func arrayAddTypeSuffixToMapKeys(rows []map[string]interface{}, mapping map[string]interface{}) {
	var excludes map[string]interface{}
	if len(mapping) > 0 {
		excludes, _ = mapping["properties"].(map[string]interface{})
	}
	for _, row := range rows {
		addTypeSuffixToMapKeys(row, excludes)
	}
}

func addTypeSuffixToMapKeys(
	m map[string]interface{}, excludes map[string]interface{},
) {
	suffixes := make(map[string]uint8)
	for k, v := range m {
		var ex interface{}
		var exMap map[string]interface{}
		if len(excludes) > 0 {
			ex = excludes[k]
			exMap, _ = ex.(map[string]interface{})
		}
		if suffix := getTypeSuffix(v, exMap); suffix > 0 && ex == nil {
			suffixes[k] = suffix
		}
	}
	for k, suffix := range suffixes {
		m[k+string([]byte{'_', suffix})] = m[k]
	}
	for k := range suffixes {
		delete(m, k)
	}
}

func getTypeSuffix(v interface{}, mapping map[string]interface{}) uint8 {
	// Unmarshal into interface, generates the 5 data type, see encoding/json/#Unmarshal
	switch value := v.(type) {
	case string:
		return 's'
	case json.Number:
		if strings.IndexByte(string(value), '.') >= 0 {
			return 'f'
		} else {
			return 'i'
		}
	case bool:
		return 'b'
	case map[string]interface{}:
		var excludes map[string]interface{}
		if len(mapping) > 0 {
			excludes, _ = mapping["properties"].(map[string]interface{})
		}
		addTypeSuffixToMapKeys(value, excludes)
		return 'o'
	case []interface{}:
		return getArrayTypeSuffix(value, mapping)
	}
	return 0
}

func getArrayTypeSuffix(slice []interface{}, mapping map[string]interface{}) uint8 {
	if len(slice) == 0 {
		return 0
	}
	suffix := getTypeSuffix(slice[0], mapping)
	if suffix == 0 {
		return 0
	}
	for i := 1; i < len(slice); i++ {
		if suffix != getTypeSuffix(slice[i], mapping) {
			return 0
		}
	}
	return suffix
}
