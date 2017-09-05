package reader

import (
	"bytes"
	"encoding/json"
	"os"
)

func (r *Reader) parseRow(line []byte) map[string]interface{} {
	if len(line) == 0 {
		return nil
	}
	var row map[string]interface{}
	if err := json.Unmarshal(line, &row); err == nil {
		return row
	} else {
		if line = bytes.TrimSpace(line); len(line) > 0 {
			r.logger.Errorf("reader: unmarshal json error(%v): %s", err, line)
		}
		return nil
	}
}

// 把文件seek到上次保存的位置
func (r *Reader) seekToSavedOffset() {
	if offset := r.offsetFile.read(); offset > 0 {
		if size := r.size(); size >= 0 && offset <= size {
			if _, err := r.file.Seek(offset, os.SEEK_SET); err == nil {
				r.logger.Printf("reader: seekToSavedOffset: %d, size: %d", offset, size)
			} else {
				r.logger.Errorf("reader: seekToSavedOffset: %d, size: %d, error: %v", offset, size, err)
			}
		} else {
			r.logger.Errorf("reader: seekToSavedOffset: offset %d exceeds size %d", offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (r *Reader) seekFrontIfTruncated() {
	if offset := r.offset(); offset > 0 {
		if size := r.size(); size >= 0 && offset > size {
			if _, err := r.file.Seek(0, os.SEEK_SET); err == nil {
				r.logger.Printf("reader: seekFront(offset: %d, size: %d)", offset, size)
				r.reader.Reset(r.file)
			} else {
				r.logger.Errorf(
					"reader: seekFront(offset: %d, size: %d) error: %v", offset, size, err,
				)
			}
		}
	}
}

func (r *Reader) offset() int64 {
	if offset, err := r.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		r.logger.Errorf("reader: get offset error: %v", err)
		return -1
	}
}

func (r *Reader) size() int64 {
	if fi, err := r.file.Stat(); err == nil {
		return fi.Size()
	} else {
		r.logger.Errorf("reader: get size error: %v", err)
		return -1
	}
}
