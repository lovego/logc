package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logFile *os.File

func init() {
	var logPath = `log/logc.log`
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		panic(err)
	}
	var err error
	if logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		panic(err)
	}
}

func writeLog(data ...string) {
	content := []string{time.Now().String()}
	content = append(content, data...)
	content = append(content, "\n")
	if _, err := logFile.WriteString(strings.Join(content, ` `)); err != nil {
		panic(err)
	}
}
