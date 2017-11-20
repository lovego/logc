package reader

import (
	"time"

	loggerpkg "github.com/lovego/xiaomei/utils/logger"
)

var batchSize = 100 * 1024
var batchWait = time.Second

type Batch struct {
	Size int    `yaml:"size"` // <= 0 means use default value
	Wait string `yaml:"wait"` // empty or < 0 means use default value, 0 means don't wait
}

func Setup(batch Batch, logger *loggerpkg.Logger) {
	if batch.Size > 0 {
		batchSize = batch.Size
	}
	if batch.Wait != `` {
		duration, err := time.ParseDuration(batch.Wait)
		if err != nil {
			logger.Fatalf("parse batch.wait error: %v", err)
		}
		if duration >= 0 {
			batchWait = duration
		}
	}
}
