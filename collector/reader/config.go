package reader

import (
	"time"

	loggerpkg "github.com/lovego/logger"
)

var (
	maxLineSize = 800 * 1024
	batchSize   = 100 * 1024
	batchWait   = time.Second
)

// The max possible size in one Read is: (batchSize - 1 + maxLineSize).

type Config struct {
	MaxLineSize int   `yaml:"maxLineSize"`
	Batch       Batch `yaml:"batch"`
}

type Batch struct {
	Size int    `yaml:"size"` // <= 0 means use default value
	Wait string `yaml:"wait"` // empty or < 0 means use default value, 0 means don't wait
}

func Setup(config Config, logger *loggerpkg.Logger) {
	if config.MaxLineSize > 0 {
		maxLineSize = config.MaxLineSize
	}
	if config.Batch.Size > 0 {
		batchSize = config.Batch.Size
	}
	if config.Batch.Wait != `` {
		duration, err := time.ParseDuration(config.Batch.Wait)
		if err != nil {
			logger.Fatalf("parse batch.wait error: %v", err)
		}
		if duration >= 0 {
			batchWait = duration
		}
	}
}
