package elastic_search

type Config struct {
	Addrs           []string `yaml:"elasticsearch"`
	Index           string   `yaml:"index"`
	timeSeriesIndex *timeSeriesIndex
	IndexKeep       int                               `yaml:"indexKeep"`
	Type            string                            `yaml:"type"`
	Mapping         map[string]map[string]interface{} `yaml:"mapping"`
	TimeField       string                            `yaml:"timeField"`
	TimeFormat      string                            `yaml:"timeFormat"`
}

func New() {
	if file.Index == `` {
		log.Fatalf("index missing for file: %s", file.Path)
	}
	var err error
	if file.TimeSeriesIndex, err = parseTimeSeriesIndex(file.Index); err != nil {
		log.Fatalf("file %s: %v", file.Path, err)
	}

	if file.Type == `` {
		log.Fatalf("type missing for file: %s", file.Path)
	}
	if file.TimeSeriesIndex != nil && file.TimeField == `` {
		log.Fatalf("timeField missing for file: %s, which has a time series index.", file.Path)
	}
	if file.TimeField != `` && file.TimeFormat == `` {
		file.TimeFormat = time.RFC3339
	}
	file.cleanMapping()
}

func (file *File) cleanMapping() {
	for _, field := range file.Mapping {
		for k, v := range field {
			if data, ok := v.(map[interface{}]interface{}); ok {
				field[k] = convertMapKeyToStr(data)
			}
		}
	}
}

func convertMapKeyToStr(data map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		if mapData, ok := v.(map[interface{}]interface{}); ok {
			result[k.(string)] = convertMapKeyToStr(mapData)
		} else {
			result[k.(string)] = v
		}
	}
	return result
}
