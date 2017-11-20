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
		case <-c.closeEvent:
			c.collect()
			c.close()
			return
		}
	}
}

func (c *Collector) collect() {
	for {
		rows, drain := c.reader.Read()
		if len(rows) > 0 {
			c.output.Write(rows)
			c.logger.Printf("%d, %s\n", len(rows), c.reader.SaveOffset())
		}
		if drain {
			break
		}
	}
}

func (c *Collector) close() {
	c.logger.Printf("collector close")
	c.reader.Close()
	if err := c.logFile.Close(); err != nil {
		logger.Errorf("logger: close error: %v", err)
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
