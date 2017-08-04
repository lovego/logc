package collector

import (
	"log"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type ifcEvent interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	Opened() bool
	Reopen()
}

type pusherIfc interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	source      ifcEvent
	pusher      pusherIfc
	logger      *log.Logger
	writeEvent  chan struct{}
	createEvent chan struct{}
}

func New(path string, source ifcEvent, pusher pusherIfc, logger *log.Logger) *Collector {
	c := &Collector{
		source: source, pusher: pusher, logger: logger,
		writeEvent:  make(chan struct{}, 1),
		createEvent: make(chan struct{}, 1),
	}
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
		case <-c.createEvent:
			if c.source.Opened() {
				utils.Protect(c.collect)
			}
			c.source.Reopen()
		case <-c.writeEvent:
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

func (c *Collector) NotifyWrite() {
	select {
	case c.writeEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) NotifyCreate() {
	select {
	case c.createEvent <- struct{}{}:
	default:
	}
}
