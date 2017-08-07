/* Package watch 监听文件变动，并通知相应的日志收集器。只关心两种变动事件：
1. Write
	通知收集器文件有写入，应该进行日志收集。
2. Create
	通知收集器文件被创建，可能是新建或移动后重新创建。应该重新打开该文件。

日志旋转（先移动再重新创建）
如果在收集日志过程中，文件被移动，因为fsnotify是基于文件名来识别文件的，就监听不到该文件的写变动。
既然做不到文件移动后继续收集， 所以我们就只监听文件所在的目录，来获取文件的写入和创建事件。
在收到创建事件时，收集器如果有当前已经打开的文件，会将当前打开的文件的所有内容收集完，
再重新打开文件以保证在日志旋转过程中没有内容丢失。

TODO
文件移动后，通知收集器改用轮询方式继续收集。文件重新创建后，切换回事件通知方式收集。
*/
package watch

import (
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Collector interface {
	NotifyWrite()
	NotifyCreate()
}

func Watch(collectors map[string]Collector) {
	files := getFiles(collectors)
	dirsWatcher := getWatcher(getDirs(files))

	defer dirsWatcher.Close()

	for {
		select {
		case err := <-dirsWatcher.Errors:
			log.Printf("dirs watcher error: %v\n", err)
		case event := <-dirsWatcher.Events:
			log.Println(event)
			if collector := collectors[strings.TrimPrefix(event.Name, `./`)]; collector != nil {
				if event.Op&fsnotify.Write == fsnotify.Write {
					collector.NotifyWrite()
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					collector.NotifyCreate()
				}
			}
		}
	}
}

func getWatcher(paths []string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("fsnotify.NewWatcher error: %v\n", err)
	}

	for _, path := range paths {
		if err := watcher.Add(path); err == nil {
			log.Printf("watch %s ", path)
		} else {
			log.Printf("watcher.Add %s error: %v\n", path, err)
		}
	}
	return watcher
}

func getFiles(collectors map[string]Collector) (files []string) {
	for file := range collectors {
		files = append(files, file)
	}
	return
}

func getDirs(files []string) (dirs []string) {
	m := make(map[string]bool)
	for _, path := range files {
		m[filepath.Dir(path)] = true
	}
	for dir := range m {
		dirs = append(dirs, dir)
	}
	return
}
