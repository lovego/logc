/* Package watch 监听文件变动，并通知相应的日志收集器。
如果在收集日志过程中，文件被移动，因为fsnotify是基于文件名来识别文件的，就监听不到该文件的写变动。
既然fsnotify做不到文件移动后继续监听， 所以我们就只监听文件所在的目录，来获取文件的写入和创建事件。
在我们检测到文件移动后，我们将收集器映射到新的文件名。
我们只关心3种变动事件：
1. Write
	通知收集器文件有写入，应该进行日志收集。
2. Create
	a. 文件被重命名：对比当前文件和所有已经打开的文件，如果是同一文件，更新收集器Map。
	b. 目标文件被创建: 创建收集器收集该文件。
3. Remove
  文件被删除，关闭收集器
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
	NotifyClose()
	OpenedSameFile(os.FileInfo) bool
	Printf(format string, v ...interface{})
}

func Watch(collectorMakers map[string]func() []Collector) {
	collectorsMap := getCollectors(collectorMakers)
	dirsWatcher := getDirsWatcher(collectorMakers)
	defer dirsWatcher.Close()

	log.Printf("\033[32mlogc started.\033[0m")

	for {
		select {
		case err := <-dirsWatcher.Errors:
			log.Printf("dirs watcher error: %v\n", err)
		case event := <-dirsWatcher.Events:
			path := strings.TrimPrefix(event.Name, `./`)
			if event.Op&fsnotify.Write == fsnotify.Write {
				for _, collector := range collectorsMap[path] {
					collector.NotifyWrite()
				}
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				if !handleRename(path, collectorsMap) {
					handleCreate(path, collectorsMap, collectorMakers)
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				handleRemove(path, collectorsMap)
			}
		}
	}
}

func handleRename(path string, collectorsMap map[string][]Collector) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Printf("watch: %v", err)
		return true
	}
	for oldPath, collectors := range collectorsMap {
		if openedSameFile(collectors, fi) {
			if oldPath != path {
				handleRemove(path, collectorsMap)
				delete(collectorsMap, oldPath)
				collectorsMap[path] = collectors
				for _, collector := range collectors {
					collector.Printf("rename from %s to %s", oldPath, path)
				}
			}
			return true
		}
	}
	return false
}

func handleCreate(
	path string, collectorsMap map[string][]Collector, collectorMakers map[string]func() []Collector,
) {
	handleRemove(path, collectorsMap)
	if maker := collectorMakers[path]; maker != nil {
		if collectors := maker(); len(collectors) > 0 {
			collectorsMap[path] = collectors
		}
	}
}

func handleRemove(path string, collectorsMap map[string][]Collector) {
	for _, collector := range collectorsMap[path] {
		collector.NotifyClose()
	}
	delete(collectorsMap, path)
}
