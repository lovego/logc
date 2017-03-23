package main

import (
	"fmt"
	"time"
)

func printLog(logs ...interface{}) {
	loc := time.FixedZone(`China`, 28800)
	now := time.Now().In(loc)
	fmt.Print(now)
	fmt.Println(logs...)
}
