package collector

import (
	"os"
)

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

func (c *Collector) NotifyWrite() {
	select {
	case c.writeEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) NotifyRemove() {
	select {
	case c.removeEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) OpenedSameFile(fi os.FileInfo) bool {
	return c.reader.SameFile(fi)
}

func (c *Collector) Printf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}
