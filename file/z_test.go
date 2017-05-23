package file

import (
	"fmt"
	"testing"
)

func TestReadOffset(t *testing.T) {
	readOffset()
	fmt.Println(offsetData.m)
}

func TestUpdateOffset(t *testing.T) {
	offsetData.m = map[string]int64{
		`/logs/api/app.log`: 0,
		`/logs/api/app.err`: 0,
		`/logs/api/web.log`: 0,
	}
	fmt.Println(updateOffset())
}
