package logger

import (
	"log"
	"os"

	"github.com/lovego/xiaomei/utils/fs"
)

type Logger struct {
	path string
	file *os.File
	*log.Logger
}

func New(path string) *Logger {
	path += `.logc`
	l := &Logger{path: path}
	if file, err := fs.OpenAppend(path); err == nil {
		l.file = file
		l.Logger = log.New(file, ``, log.LstdFlags)
		return l
	} else {
		log.Printf("logger: open %s error: %v", path, err)
		return nil
	}
}

func (l *Logger) Get() *log.Logger {
	return l.Logger
}

func (l *Logger) Rename(newPath string) {
	newPath += `.logc`
	if err := os.Rename(l.path, newPath); err == nil {
		l.path = newPath
	} else {
		l.Printf("logger: rename %s to %s error: %v", l.path, newPath, err)
	}
}

func (l *Logger) Remove() {
	if err := l.file.Close(); err != nil {
		l.Printf("logger: close %s error: %v", l.path, err)
	}
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		l.Printf("logger: remove %s error: %v", l.path, err)
	}
}
