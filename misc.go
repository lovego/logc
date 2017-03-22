package main

import (
	"fmt"
	"time"
)

func printLog(logs ...interface{}) {
	loc := time.FixedZone(`China`, 28800)
	now := time.Now().In(loc)
	fmt.Println(append([]interface{}{now}, logs...)...)
}
