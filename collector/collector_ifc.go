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

func (c *Collector) NotifyRename() {
	select {
	case c.renameEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) NotifyDestroy() {
	select {
	case c.destroyEvent <- struct{}{}:
	default:
	}
}

func (c *Collector) OpenedSameFile(os.FileInfo) bool {
	return false
}
