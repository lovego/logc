package collector

import (
	"os"
)

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
