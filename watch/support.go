package watch

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

func getCollectors(collectorMakers map[string]func() []Collector) map[string][]Collector {
	collectorsMap := make(map[string][]Collector)
	for path, maker := range collectorMakers {
		if collectors := maker(); len(collectors) > 0 {
			collectorsMap[path] = collectors
		}
	}
	return collectorsMap
}

func getDirsWatcher(collectorMakers map[string]func() []Collector) *fsnotify.Watcher {
	dirs := make(map[string]bool)
	for path := range collectorMakers {
		dirs[filepath.Dir(path)] = true
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("fsnotify.NewWatcher error: %v\n", err)
	}

	for dir := range dirs {
		if err := watcher.Add(dir); err == nil {
			log.Printf("watch %s ", dir)
		} else {
			log.Printf("watcher.Add %s error: %v\n", dir, err)
		}
	}
	return watcher
}

func openedSameFile(collectors []Collector, fi os.FileInfo) bool {
	for _, collector := range collectors {
		if collector.OpenedSameFile(fi) {
			return true
		}
	}
	return false
}
