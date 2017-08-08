/* Package watch 监听文件变动，并通知相应的日志收集器。
如果在收集日志过程中，文件被移动，因为fsnotify是基于文件名来识别文件的，就监听不到该文件的写变动。
既然fsnotify做不到文件移动后继续监听， 所以我们就只监听文件所在的目录，来获取文件的写入和创建事件。
在我们检测到文件移动后，我们将收集器映射到新的文件名。
我们只关心两种变动事件：
1. Write
	通知收集器文件有写入，应该进行日志收集。
2. Create
	a. 文件被重命名：对比当前文件和所有已经打开的文件，如果是同一文件，更新收集器及收集器Map。
	b. 目标文件被创建: 应该通知收集器重新打开该文件。
3. Remove
  文件被删除，销毁收集器
*/
package watch

import (
	"log"
	"os"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Collector interface {
	NotifyWrite()
	NotifyRename(newPath string)
	OpenedSameFile(os.FileInfo) bool
	Destroy()
}

func Watch(collectorMakers map[string]func() Collector) {
	collectors := getCollectors(collectorMakers)
	dirsWatcher := getWatcher(getDirs(collectors))
	defer dirsWatcher.Close()

	for {
		select {
		case err := <-dirsWatcher.Errors:
			log.Printf("dirs watcher error: %v\n", err)
		case event := <-dirsWatcher.Events:
			log.Println(event)

			path := strings.TrimPrefix(event.Name, `./`)
			if event.Op&fsnotify.Write == fsnotify.Write {
				if collector := collectors[path]; collector != nil {
					collector.NotifyWrite()
				}
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				if !handleRename(path, collectors) {
					handleCreate(path, collectors, collectorMakers)
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				handleRemove(path, collectors)
			}
		}
	}
}

func handleRename(path string, collectors map[string]Collector) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Printf("stat %s error: %v", path, err)
		return true
	}
	for oldPath, collector := range collectors {
		if collector.OpenedSameFile(fi) {
			if oldPath != path {
				if coll := collectors[path]; coll != nil {
					coll.Destroy()
				}
				collector.NotifyRename(path)
				delete(collectors, oldPath)
				collectors[path] = collector
			}
			return true
		}
	}
	return false
}

func handleCreate(
	path string, collectors map[string]Collector, collectorMakers map[string]func() Collector,
) {
	handleRemove(path, collectors)
	if maker := collectorMakers[path]; maker != nil {
		collectors[path] = maker()
	}
}

func handleRemove(path string, collectors map[string]Collector) {
	if collector := collectors[path]; collector != nil {
		collector.Destroy()
		delete(collectors, path)
	}
}
