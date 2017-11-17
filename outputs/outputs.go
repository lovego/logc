package pusher

import (
	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/xiaomei/utils/logger"
)

type Getter struct {
	file config.File
}

func NewGetter(file config.File) collector.PusherGetter {
	return &Getter{file: file}
}

func (g *Getter) Get(log *logger.Logger) collector.Pusher {
	return &Pusher{file: g.file, logger: log}
}
