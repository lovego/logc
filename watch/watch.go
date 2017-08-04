package watch

import (
	"log"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type Collector interface {
	Notify()
	Reopen()
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
				if event.Op&fsnotify.Create > 0 {
					collector.Reopen()
				}
				if event.Op&fsnotify.Write > 0 {
					collector.Notify()
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
			if relPath, err := filepath.Rel(dir, path); err != nil {
				log.Panic(err)
			} else if path != relPath {
				delete(collectors, path)
				collectors[relPath] = collector
			}
		}
	}
	for dir := range m {
		dirs = append(dirs, dir)
	}
	return
}
