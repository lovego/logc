package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logPath = `log/logc.log`
var f *os.File

func init() {
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		panic(err)
	}
	if f == nil {
		var err error
		if f, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
			panic(err)
		}
	}
}

func printLog(logs ...interface{}) {
	loc := time.FixedZone(`China`, 28800)
	now := time.Now().In(loc)
	fmt.Print(now)
	fmt.Println(logs...)
}

func writeLog(data ...string) {
	now := time.Now().In(time.FixedZone(`China`, 28800))
	content := []string{now.String()}
	content = append(content, data...)
	content = append(content, "\n")
	if _, err := f.WriteString(strings.Join(content, ` `)); err != nil {
		panic(err)
	}
}
