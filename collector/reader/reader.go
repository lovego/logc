package reader

import (
	"bufio"
	"io"
	"os"

	"github.com/lovego/xiaomei/utils/logger"
)

type Reader struct {
	file       *os.File
	offsetFile *offsetFile
	buffered   []byte
	reader     *bufio.Reader
	logger     *logger.Logger
}

func New(file *os.File, offsetPath string, logger *logger.Logger) *Reader {
	r := &Reader{file: file, reader: bufio.NewReader(file), logger: logger}
	if offsetFile := newOffsetFile(offsetPath, logger); offsetFile != nil {
		r.offsetFile = offsetFile
	} else {
		return nil
	}
	r.seekToSavedOffset()

	return r
}

var batchSize = 100 * 1024

func SetBatchSize(size int) {
	if size > 0 {
		batchSize = size
	}
}

func (r *Reader) Read() (rows []map[string]interface{}, drain bool) {
	for size := 0; size < batchSize; {
		line, err := r.readLine()
		if len(line) > 0 {
			if row := r.parseRow(line); row != nil {
				rows = append(rows, row)
				size += len(line)
			}
		}
		if err != nil {
			if err == io.EOF {
				if size == 0 && len(line) == 0 {
					r.seekFrontIfTruncated()
				}
			} else {
				r.logger.Errorf("reader: read error: %v", err)
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
		r.logger.Errorf("reader: stat error: %v", err)
		return false
	} else {
		return os.SameFile(thisFi, fi)
	}
}

func (r *Reader) Close() {
	if err := r.file.Close(); err != nil {
		r.logger.Errorf("reader: close error: %v", err)
	}
	r.offsetFile.Close()
}
