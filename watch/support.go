package watch

import (
	"log"
	"path/filepath"

	// "github.com/lovego/xiaomei/utils/fs"
	"gopkg.in/fsnotify.v1"
)

func getCollectors(collectorMakers map[string]func() Collector) map[string]Collector {
	collectors := make(map[string]Collector)
	for path, maker := range collectorMakers {
		collectors[path] = maker()
	}
	return collectors
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

func getDirs(collectors map[string]Collector) map[string]bool {
	m := make(map[string]bool)
	for path := range collectors {
		m[filepath.Dir(path)] = true
	}
	return m
}
