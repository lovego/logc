package collector

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/lovego/logc/collector/logger"
	"github.com/lovego/logc/collector/reader"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Remove()
}

type Reader interface {
	Read() (rows []map[string]interface{}, drain bool)
	SaveOffset() string
	SameFile(os.FileInfo) bool
	Remove()
}

type PusherGetter interface {
	Get(*log.Logger) Pusher
}

type Pusher interface {
	Push(rows []map[string]interface{})
}

type Collector struct {
	logger      Logger
	reader      Reader
	pusher      Pusher
	writeEvent  chan struct{}
	removeEvent chan struct{}
}

func New(path string, pusherGetter PusherGetter) *Collector {
	var f *os.File
	var l *logger.Logger
	if f = openFile(path); f != nil {
		if logcPath := getLogcPath(path, f); logcPath != `` {
			if l = logger.New(logcPath + `.log`); l != nil {
				l.Printf("collect %s", path)
				if r := reader.New(f, logcPath+`.offset`, l.Get()); r != nil {
					c := &Collector{
						logger: l, reader: r, pusher: pusherGetter.Get(l.Get()),
						writeEvent:  make(chan struct{}, 1),
						removeEvent: make(chan struct{}, 1),
					}
					go c.loop()
					return c
				}
			}
		}
	}
	freeResource(f, l, path)
	return nil
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
