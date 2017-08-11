package reader

import (
	"bufio"
	"io"
	"log"
	"os"
)

type Reader struct {
	file       *os.File
	offsetFile *offsetFile
	reader     *bufio.Reader
	logger     *log.Logger
}

func New(file *os.File, offsetPath string, logger *log.Logger) *Reader {
	r := &Reader{file: file, reader: bufio.NewReader(file), logger: logger}
	if offsetFile := newOffsetFile(offsetPath, logger); offsetFile != nil {
		r.offsetFile = offsetFile
	} else {
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
				r.logger.Printf("reader: read error: %v", err)
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
		r.logger.Printf("reader: stat error: %v", err)
		return false
	} else {
		return os.SameFile(thisFi, fi)
	}
}

func (r *Reader) Close() {
	if err := r.file.Close(); err != nil {
		r.logger.Printf("reader: close error: %v", err)
	}
	r.offsetFile.Close()
}
