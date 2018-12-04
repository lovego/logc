package elasticsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

func ExampleAddTypeSuffixToMapKeys_empty() {
	m := map[string]interface{}{}
	addTypeSuffixToMapKeys(m, nil)
	fmt.Println(m)

	// Output:
	// map[]
}

func ExampleAddTypeSuffixToMapKeys_basic() {
	m := map[string]interface{}{
		"bool":     true,
		"float":    json.Number("12.3"),
		"int":      json.Number("123"),
		"intArray": []interface{}{json.Number("123")},
		"str":      "str",
	}
	addTypeSuffixToMapKeys(m, nil)
	sortPrint(m)

	// Output:
	// bool_b: true
	// float_f: 12.3
	// intArray_i: [123]
	// int_i: 123
	// str_s: str
}

func ExampleAddTypeSuffixToMapKeys_basicWithExcludes() {
	m := map[string]interface{}{
		"bool":     true,
		"float":    json.Number("12.3"),
		"int":      json.Number("123"),
		"intArray": []interface{}{json.Number("123")},
		"str":      "str",
	}
	addTypeSuffixToMapKeys(m, map[string]interface{}{
		"bool":  "bool",
		"float": map[string]string{"type": "float"},
	})
	sortPrint(m)

	// Output:
	// bool: true
	// float: 12.3
	// intArray_i: [123]
	// int_i: 123
	// str_s: str
}

func ExampleAddTypeSuffixToMapKeys_object() {
	m := map[string]interface{}{
		"object": map[string]interface{}{
			"bool":      true,
			"boolArray": []interface{}{true, false},
			"float":     json.Number("12.3"),
			"int":       json.Number("123"),
			"str":       "str",
		},
	}
	addTypeSuffixToMapKeys(m, nil)
	sortPrint(m["object_o"].(map[string]interface{}))
	// Output:
	// boolArray_b: [true false]
	// bool_b: true
	// float_f: 12.3
	// int_i: 123
	// str_s: str
}

func ExampleAddTypeSuffixToMapKeys_objectWithExcludes() {
	m := map[string]interface{}{
		"object": map[string]interface{}{
			"bool":      true,
			"boolArray": []interface{}{true, false},
			"float":     json.Number("12.3"),
			"int":       json.Number("123"),
			"str":       "str",
		},
	}
	addTypeSuffixToMapKeys(m, map[string]interface{}{
		"object": map[string]interface{}{
			"properties": map[string]interface{}{
				"bool":  "bool",
				"float": map[string]string{"type": "float"},
			},
		},
	})
	sortPrint(m["object"].(map[string]interface{}))
	// Output:
	// bool: true
	// boolArray_b: [true false]
	// float: 12.3
	// int_i: 123
	// str_s: str
}

func ExampleAddTypeSuffixToMapKeys_objectWithExcludes2() {
	m := map[string]interface{}{
		"object": map[string]interface{}{
			"bool":      true,
			"boolArray": []interface{}{true, false},
			"float":     json.Number("12.3"),
			"int":       json.Number("123"),
			"str":       "str",
		},
	}
	addTypeSuffixToMapKeys(m, map[string]interface{}{
		"object": map[string]interface{}{},
	})
	sortPrint(m["object"].(map[string]interface{}))
	// Output:
	// boolArray_b: [true false]
	// bool_b: true
	// float_f: 12.3
	// int_i: 123
	// str_s: str
}

func ExampleUnmarshal() {
	var data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(`{
    "int": 1,
    "intArray": [1,2],
    "objectArray": [ { "k": "v" } ],
    "map": { "k": "v" }
  }`), &data); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%T\n", data["int"])
	fmt.Printf("%T\n", data["intArray"])
	fmt.Printf("%T\n", data["objectArray"])
	fmt.Printf("%T\n", data["map"])
	// Output:
	// float64
	// []interface {}
	// []interface {}
	// map[string]interface {}
}

func sortPrint(m map[string]interface{}) {
	var slice [][2]interface{}
	for k, v := range m {
		slice = append(slice, [2]interface{}{k, v})
	}
	sort.Slice(slice, func(i, j int) bool { return slice[i][0].(string) < slice[j][0].(string) })
	for _, row := range slice {
		fmt.Printf("%s: %v\n", row[0], row[1])
	}
}
