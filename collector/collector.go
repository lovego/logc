package collector

import (
	"os"
)

type readerIfc interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Rename(string)
	Remove()
}

type loggerIfc interface {
	Printf(format string, v ...interface{})
	Rename(string)
	Remove()
}

type pusherIfc interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	path        string
	reader      readerIfc
	logger      loggerIfc
	pusher      pusherIfc
	writeEvent  chan struct{}
	renameEvent chan string
	removeEvent chan struct{}
}

func New(path string, reader readerIfc, logger loggerIfc, pusher pusherIfc) *Collector {
	c := &Collector{
		path: path, reader: reader, pusher: pusher, logger: logger,
		writeEvent:  make(chan struct{}, 1),
		renameEvent: make(chan string, 1),
		removeEvent: make(chan struct{}, 1),
	}

	go c.loop()
	return c
}

func (c *Collector) loop() {
	c.logger.Printf("collect %s", c.path)
	c.collect() // collect existing data.
	for {
		select {
		case newPath := <-c.renameEvent:
			c.rename(newPath)
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

func (c *Collector) rename(newPath string) {
	c.logger.Printf("rename from %s to %s", c.path, newPath)
	c.path = newPath
	c.reader.Rename(newPath)
	c.logger.Rename(newPath)
}

func (c *Collector) remove() {
	c.logger.Printf("removed %s", c.path)
	c.reader.Remove()
	c.logger.Remove()
}
