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
	filesWatcher := getWatcher(files)
	dirsWatcher := getWatcher(getDirs(files))

	defer filesWatcher.Close()
	defer dirsWatcher.Close()

	for {
		select {
		case err := <-filesWatcher.Errors:
			log.Printf("files watcher error: %v\n", err)
		case err := <-dirsWatcher.Errors:
			log.Printf("dirs watcher error: %v\n", err)
		case event := <-filesWatcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				if collector := collectors[event.Name]; collector != nil {
					log.Println(event)
					collector.NotifyWrite()
				}
			}
		case event := <-dirsWatcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				if collector := collectors[strings.TrimPrefix(event.Name, `./`)]; collector != nil {
					log.Println(event)
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
		if err := watcher.Add(path); err != nil {
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
