package reader

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/lovego/logger"
)

type Reader struct {
	file        *os.File
	offsetFile  *offsetFile
	buffered    []byte
	reader      *bufio.Reader
	logger      *logger.Logger
	collectorId string
}

func New(collectorId string, file *os.File, offsetPath string, logger *logger.Logger) *Reader {
	r := &Reader{
		collectorId: collectorId, file: file, reader: bufio.NewReader(file), logger: logger,
	}
	if offsetFile := newOffsetFile(offsetPath, logger); offsetFile != nil {
		r.offsetFile = offsetFile
	} else {
		return nil
	}
	r.seekToSavedOffset()

	return r
}

// drain 是否读完了所有内容
func (r *Reader) Read() ([]map[string]interface{}, bool) {
	rows, drain, size := r.readSize(batchSize)
	if len(rows) > 0 && drain && batchWait > 0 {
		time.Sleep(batchWait)
		if rows2, drain2, _ := r.readSize(batchSize - size); len(rows2) > 0 {
			rows = append(rows, rows2...)
			drain = drain2
		}
	}
	if len(rows) == 0 {
		r.seekFrontIfTruncated()
	}
	return rows, drain
}

func (r *Reader) readSize(targetSize int) (rows []map[string]interface{}, drain bool, size int) {
	// 读到文件末尾时返回错误：io.EOF
	var err error
	for err == nil && size < targetSize {
		var line []byte
		if line, err = r.readLine(); len(line) > 0 {
			if row := r.parseLine(line); len(row) > 0 {
				rows = append(rows, row)
				size += len(line)
			}
		}
	}
	if err != nil {
		drain = true // 读到文件末尾或者读取出错都认为都完了所有内容
		if err != io.EOF {
			r.logger.Errorf("%s: reader: read error: %v", r.collectorId, err)
		}
	}
	return
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
		r.logger.Errorf("%s: reader: stat error: %v", r.collectorId, err)
		return false
	} else {
		return os.SameFile(thisFi, fi)
	}
}

func (r *Reader) Close() {
	if err := r.file.Close(); err != nil {
		r.logger.Errorf("%s: reader: close error: %v", r.collectorId, err)
	}
	r.offsetFile.Close()
}
