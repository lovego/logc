package reader

import (
	"bufio"
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

func (r *Reader) SameFile(fi os.FileInfo) bool {
	if thisFi, err := r.file.Stat(); err != nil {
		r.logger.Printf("reader: file.Stat %s error: %v", r.path, err)
		return false
	} else {
		return os.SameFile(thisFi, fi)
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
