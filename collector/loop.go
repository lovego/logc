package collector

import (
	"os"
)

func (c *Collector) loop() {
	// collect existing data.
	if !c.collect() {
		c.logger.Errorf("%s: collector exited.", c.id)
		c.close()
		return
	}
	for {
		select {
		case <-c.writeEvent:
			if !c.collect() {
				c.logger.Errorf("%s: collector exited.", c.id)
				c.close()
				return
			}
		case <-c.closeEvent:
			c.collect()
			c.logger.Printf("collector close")
			c.close()
			return
		}
	}
}

func (c *Collector) collect() bool {
	for {
		rows, drain := c.reader.Read()
		if len(rows) > 0 {
			if c.output.Write(rows) {
				c.logger.Printf("%d, %s\n", len(rows), c.reader.SaveOffset())
			} else {
				return false
			}
		}
		if drain {
			break
		}
	}
	return true
}

func (c *Collector) close() {
	c.reader.Close()
	if err := c.logFile.Close(); err != nil {
		c.logger.Errorf("%s: close logger error: %v.", c.id, err)
	}
}

func (c *Collector) NotifyWrite() {
	select {
	case c.writeEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) NotifyClose() {
	select {
	case c.closeEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) OpenedSameFile(fi os.FileInfo) bool {
	return c.reader.SameFile(fi)
}

func (c *Collector) Printf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}
