package collector

import (
	"os"
)

type readerIfc interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Remove()
}

type loggerIfc interface {
	Printf(format string, v ...interface{})
	Remove()
}

type pusherIfc interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	reader      readerIfc
	logger      loggerIfc
	pusher      pusherIfc
	writeEvent  chan struct{}
	removeEvent chan struct{}
}

func New(path string, reader readerIfc, logger loggerIfc, pusher pusherIfc) *Collector {
	c := &Collector{
		reader: reader, pusher: pusher, logger: logger,
		writeEvent:  make(chan struct{}, 1),
		removeEvent: make(chan struct{}, 1),
	}

	c.logger.Printf("collect %s", path)
	go c.loop()
	return c
}

func (c *Collector) loop() {
	c.collect() // collect existing data.
	for {
		select {
		case <-c.writeEvent:
			c.collect()
		case <-c.removeEvent:
			c.collect()
			c.remove()
			return
		}
	}
}

func (c *Collector) collect() {
	for {
		rows, drain := c.reader.Read()
		if len(rows) > 0 {
			c.pusher.Push(rows)
			c.logger.Printf("%d, %s\n", len(rows), c.reader.SaveOffset())
		}
		if drain {
			break
		}
	}
}

func (c *Collector) remove() {
	c.logger.Printf("remove")
	c.reader.Remove()
	c.logger.Remove()
}
