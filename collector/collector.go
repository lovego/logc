package collector

import (
	"log"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type sourceIfc interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	Reopen()
}

type pusherIfc interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	source       sourceIfc
	pusher       pusherIfc
	logger       *log.Logger
	sourceWrite  chan struct{}
	sourceChange chan struct{}
}

func New(path string, source sourceIfc, pusher pusherIfc, logger *log.Logger) *Collector {
	c := &Collector{source: source, pusher: pusher, logger: logger}
	log.Println(`listen ` + path)

	go c.loop(path)
	return c
}

func (c *Collector) loop(path string) {
	c.logger.Println(`listen ` + path)
	if fs.Exist(path) {
		c.collect() // collect existing data.
	}
	for {
		select {
		case <-c.sourceChange:
			utils.Protect(c.collect) // collect remaining data.
			c.source.Reopen()
		case <-c.sourceWrite:
			utils.Protect(c.collect)
		}
	}
}

func (c *Collector) collect() {
	for {
		rows, drain := c.source.Read()
		if len(rows) > 0 {
			c.pusher.Push(rows)
			c.logger.Printf("%d, %s\n", len(rows), c.source.SaveOffset())
		}
		if drain {
			break
		}
	}
}

func (c *Collector) NotifySourceWrite() {
	select {
	case c.sourceWrite <- struct{}{}:
	default:
	}
}

func (c *Collector) NotifySourceChange() {
	select {
	case c.sourceChange <- struct{}{}:
	default:
	}
}
