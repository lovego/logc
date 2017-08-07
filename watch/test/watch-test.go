package main

import (
	"log"

	"gopkg.in/fsnotify.v1"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("fsnotify.NewWatcher error: %v\n", err)
	}
	defer watcher.Close()

	for _, path := range []string{`./test.txt`} {
		if err := watcher.Add(path); err != nil {
			log.Printf("watcher.Add %s error: %v\n", path, err)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println(event)
		case err := <-watcher.Errors:
			log.Printf("notify watcher error: %v\n", err)
		}
	}
}
