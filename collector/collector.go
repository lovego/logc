package collector

import (
	"log"
	"os"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
)

type sourceIfc interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	RenameOffset(string)
	Destroy()
}

type loggerIfc interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Rename(string)
	Destroy()
}

type pusherIfc interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	source       sourceIfc
	logger       loggerIfc
	pusher       pusherIfc
	writeEvent   chan struct{}
	renameEvent  chan struct{}
	destroyEvent chan struct{}
}

func New(path string, source sourceIfc, logger loggerIfc, pusher pusherIfc) *Collector {
	c := &Collector{
		source: source, pusher: pusher, logger: logger,
		writeEvent:   make(chan struct{}, 1),
		renameEvent:  make(chan struct{}, 1),
		destroyEvent: make(chan struct{}, 1),
	}

	c.logger.Println(`listen ` + path)
	go c.loop(path)
	return c
}

func (c *Collector) loop(path string) {
	c.collect() // collect existing data.
	for {
		select {
		case <-c.renameEvent:
			utils.Protect(c.rename)
		case <-c.writeEvent:
			utils.Protect(c.collect)
		case <-c.destroyEvent:
			utils.Protect(c.destroy)
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

func (c *Collector) rename(newPath string) {
	s.source.RenameOffset(newPath)
	s.logger.Rename()
}

func (c *Collector) destroy() {
	s.source.Destroy()
	s.logger.Destroy()
}
