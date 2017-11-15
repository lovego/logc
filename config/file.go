package config

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type File struct {
	Path            string `yaml:"path"`
	Index           string `yaml:"index"`
	indexPrefix     string
	indexTimeLayout string
	indexSuffix     string
	IndexKeep       int                               `yaml:"indexKeep"`
	Type            string                            `yaml:"type"`
	Mapping         map[string]map[string]interface{} `yaml:"mapping"`
	TimeField       string                            `yaml:"timeField"`
	TimeFormat      string                            `yaml:"timeFormat"`
}

func (file *File) IsTimeSeriesIndex() bool {
	return file.indexTimeLayout != ``
}

func (file *File) GetIndex(t time.Time) string {
	if file.indexTimeLayout != `` {
		return file.indexPrefix + t.Format(file.indexTimeLayout) + file.indexSuffix
	}
	return file.Index
}

func (file *File) check() {
	if file.Path == `` {
		log.Fatalf("path missing for file: %+v", file)
	} else {
		file.Path = filepath.Clean(file.Path)
	}
	file.checkIndex()
	if file.Type == `` {
		log.Fatalf("type missing for file: %s", file.Path)
	}
	if file.TimeField != `` && file.TimeFormat == `` {
		file.TimeFormat = time.RFC3339
	}
	file.cleanMapping()
}

var timeSeriesIndexRegexp = regexp.MustCompile(`^([^<>]*)(<[^<>]+>)([^<>]*)$`)

func (file *File) checkIndex() {
	if file.Index == `` {
		log.Fatalf("index missing for file: %s", file.Path)
	}
	if strings.IndexByte(file.Index, '<') < 0 && strings.IndexByte(file.Index, '>') < 0 {
		return
	}
	if m := timeSeriesIndexRegexp.FindStringSubmatch(file.Index); len(m) == 4 {
		file.indexPrefix, file.indexTimeLayout, file.indexSuffix = m[1], m[2], m[3]
	} else {
		log.Fatalf("file %s invalid index: %s", file.Path, file.Index)
	}

	if file.indexTimeLayout != `` && file.IndexKeep == 0 {
		file.IndexKeep = 3
	}
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
