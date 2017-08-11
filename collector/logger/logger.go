package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lovego/xiaomei/utils/fs"
)

type Logger struct {
	path string
	file *os.File
	*log.Logger
}

func New(path string) *Logger {
	l := &Logger{path: path}
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		log.Printf("logger: %v", err) // os.PathError is enough
		return nil
	}
	if file, err := fs.OpenAppend(path); err == nil {
		l.file = file
		l.Logger = log.New(file, ``, log.LstdFlags)
		return l
	} else {
		log.Printf("logger: %v", err) // os.PathError is enough
		return nil
	}
}

func (l *Logger) Get() *log.Logger {
	return l.Logger
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Printf("logger: close %s error: %v", l.path, err)
	}
}
