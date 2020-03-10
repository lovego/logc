package collector

import (
	"os"
)

func (c *Collector) loop() {
	// 收集已经存在的数据
	if !c.collect() {
		c.exitActively() // 主动退出
		return
	}
	for {
		select {
		case <-c.writeEvent:
			if !c.collect() {
				c.exitActively() // 主动退出
				return
			}
		case <-c.closeEvent:
			c.collect()
			// watch通知退出
			c.logger.Infof("collector closed.")
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
				c.logger.Infof("%d, %s\n", len(rows), c.reader.SaveOffset())
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

// 主动退出，watch不知道collector已经退出，仍然会发送通知。
func (c *Collector) exitActively() {
	c.logger.Errorf("%s: collector exited.", c.id)
	c.close()
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
	c.logger.Infof(format, v...)
}
