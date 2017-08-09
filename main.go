package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/lovego/logc/collector"
	"github.com/lovego/logc/config"
	"github.com/lovego/logc/logger"
	"github.com/lovego/logc/pusher"
	"github.com/lovego/logc/reader"
	"github.com/lovego/logc/watch"
)

func main() {
	conf := config.Get()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	pusher.CreateMappings(conf.LogdAddr, conf.Files)

	collectors := make(map[string]func() watch.Collector)
	for _, file := range conf.Files {
		collectors[file.Path] = getCollectorMaker(file, conf.LogdAddr, conf.MergeJson)
	}
	watch.Watch(collectors)
}

func getCollectorMaker(file *config.File, logdAddr, mergeJson string) func() watch.Collector {
	return func() watch.Collector {
		var f *os.File
		var l *logger.Logger

		if f = openFile(file.Path); f != nil {
			if logcPath := getLogcPath(file.Path, f); logcPath != `` {
				if l = logger.New(logcPath + `.log`); l != nil {
					if r := reader.New(f, logcPath+`.offset`, l.Get()); r != nil {
						return collector.New(
							file.Path, r, l, pusher.New(logdAddr, file.Org, file.Name, mergeJson, l.Get()),
						)
					}
				}
			}
		}
		freeResource(f, l, file.Path)
		return nil
	}
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println(err) // os.PathError is enough
		}
		return nil
	}
	return file
}

func getLogcPath(path string, f *os.File) string {
	fi, err := f.Stat()
	if err != nil {
		log.Println("stat:", err)
		return ``
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || sys == nil {
		log.Printf("unexpected FileInfo.Sys(): %#v", fi.Sys())
		return ``
	}
	return filepath.Join(filepath.Dir(path), `logc`, strconv.FormatUint(sys.Ino, 10))
}

func freeResource(f *os.File, l *logger.Logger, path string) {
	if f != nil {
		if err := f.Close(); err != nil {
			log.Printf("close %s error: %v", path, err)
		}
	}
	if l != nil {
		l.Close()
	}
}
