package watch

import (
	"log"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type Collector interface {
	NotifyWrite()
	NotifyCreate()
}

func Watch(collectors map[string]Collector) {
	watcher := initWatcher(collectors)
	if watcher == nil {
		return
	}
	defer watcher.Close()

	for {
		select {
		case event := <-watcher.Events:
			collector := collectors[event.Name]
			if collector != nil {
				log.Println(event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					collector.NotifyCreate()
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					collector.NotifyWrite()
				}
			}
		case err := <-watcher.Errors:
			log.Printf("notify watcher error: %v\n", err)
		}
	}
}

func initWatcher(collectors map[string]Collector) *fsnotify.Watcher {
	dirs := getWatchDirs(collectors)
	if len(dirs) == 0 {
		return nil
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("fsnotify.NewWatcher error: %v\n", err)
	}

	hasSuccess := false
	for _, dir := range dirs {
		if err := watcher.Add(dir); err != nil {
			log.Printf("watcher.Add %s error: %v\n", dir, err)
		} else {
			hasSuccess = true
		}
	}
	if hasSuccess {
		return watcher
	} else {
		watcher.Close()
		return nil
	}
}

func getWatchDirs(collectors map[string]Collector) (dirs []string) {
	m := make(map[string]bool)
	for path, collector := range collectors {
		dir := filepath.Dir(path)
		m[dir] = true
		if dir == `.` {
			delete(collectors, path)
			collectors[`./`+path] = collector
		}
	}
	for dir := range m {
		dirs = append(dirs, dir)
	}
	return
}
