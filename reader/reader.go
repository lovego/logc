package reader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

type Reader struct {
	path       string
	file       *os.File
	offsetFile *offsetFile
	reader     *bufio.Reader
	logger     *log.Logger
}

func New(path string, logger *log.Logger) *Reader {
	file, err := os.Open(path)
	if err != nil {
		logger.Printf("reader: open %s error: %v", path, err)
		return nil
	}
	r := &Reader{path: path, file: file, reader: bufio.NewReader(file), logger: logger}
	if offsetFile := newOffsetFile(path, logger); offsetFile != nil {
		r.offsetFile = offsetFile
	} else {
		file.Close()
		return nil
	}
	r.seekToSavedOffset()

	return r
}

func (r *Reader) Read() (rows []map[string]interface{}, drain bool) {
	for size := 0; size < 2*1024*1024; {
		line, err := r.reader.ReadBytes('\n')
		if row := r.parseRow(line); row != nil {
			rows = append(rows, row)
			size += len(line)
		}
		if err != nil {
			if err == io.EOF {
				if size == 0 && len(line) == 0 {
					r.seekFrontIfTruncated()
				}
			} else {
				r.logger.Printf("reader: read %s error: %v", r.path, err)
			}
			return rows, true
		}
	}
	return rows, false
}

func (r *Reader) SaveOffset() string {
	if offset := r.offset(); offset > 0 {
		return r.offsetFile.save(offset)
	} else {
		return ``
	}
}

func (r *Reader) Rename(newPath string) {
	r.path = newPath
	r.offsetFile.rename(newPath)
}

func (r *Reader) Remove() {
	if err := r.file.Close(); err != nil {
		r.logger.Printf("reader: close %s error: %v", r.path, err)
	}
	r.offsetFile.remove()
}

func (r *Reader) parseRow(line []byte) map[string]interface{} {
	if len(line) == 0 {
		return nil
	}
	var row map[string]interface{}
	if err := json.Unmarshal(line, &row); err == nil {
		return row
	} else {
		if line = bytes.TrimSpace(line); len(line) > 0 {
			r.logger.Printf("reader: %s unmarshal json error(%v): %s", r.path, err, line)
		}
		return nil
	}
}

// 把文件seek到上次保存的位置
func (r *Reader) seekToSavedOffset() {
	if offset := r.offsetFile.read(); offset > 0 {
		if size := r.size(); size >= 0 && offset <= size {
			if _, err := r.file.Seek(offset, os.SEEK_SET); err == nil {
				r.logger.Printf("reader: %s seek to: %d, size: %d", r.path, offset, size)
			} else {
				r.logger.Printf("reader: %s seek to: %d, size: %d, error: %v", r.path, offset, size, err)
			}
		} else {
			r.logger.Printf("reader: %s offset %d exceeds size %d", r.path, offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (r *Reader) seekFrontIfTruncated() {
	if offset := r.offset(); offset > 0 {
		if size := r.size(); size >= 0 && offset > size {
			if _, err := r.file.Seek(0, os.SEEK_SET); err == nil {
				r.logger.Printf("reader: %s seekFront(offset: %d, size: %d)", r.path, offset, size)
				r.reader.Reset(r.file)
			} else {
				r.logger.Printf(
					"reader: %s seekFront(offset: %d, size: %d) error: %v", r.path, offset, size, err,
				)
			}
		}
	}
}

func (r *Reader) offset() int64 {
	if offset, err := r.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		r.logger.Printf("reader: %s get offset error: %v", r.path, err)
		return -1
	}
}

func (r *Reader) size() int64 {
	if fi, err := r.file.Stat(); err == nil {
		return fi.Size()
	} else {
		r.logger.Printf("reader: %s get size error: %v", r.path, err)
		return -1
	}
}
