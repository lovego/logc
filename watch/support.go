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

func getWatcher(paths map[string]bool) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("fsnotify.NewWatcher error: %v\n", err)
	}

	for path := range paths {
		if err := watcher.Add(path); err == nil {
			log.Printf("watch %s ", path)
		} else {
			log.Printf("watcher.Add %s error: %v\n", path, err)
		}
	}
	return watcher
}

func getDirs(collectorsMap map[string][]Collector) map[string]bool {
	m := make(map[string]bool)
	for path := range collectorsMap {
		m[filepath.Dir(path)] = true
	}
	return m
}

func openedSameFile(collectors []Collector, fi os.FileInfo) bool {
	for _, collector := range collectors {
		if collector.OpenedSameFile(fi) {
			return true
		}
	}
	return false
}
