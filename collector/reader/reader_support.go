package reader

import (
	"bytes"
	"encoding/json"
	"os"
)

func (r *Reader) readLine() ([]byte, error) {
	line, err := r.reader.ReadBytes('\n')
	if len(line) == 0 {
		return line, err
	}
	if line[len(line)-1] != '\n' {
		r.buffered = append(r.buffered, line...)
		return nil, err
	}
	if len(r.buffered) > 0 {
		line = append(r.buffered, line...)
		r.buffered = nil
	}
	return line, err
}

func (r *Reader) parseLine(line []byte) map[string]interface{} {
	var row map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(line))
	decoder.UseNumber()
	if err := decoder.Decode(&row); err == nil {
		return row
	} else {
		if line = bytes.TrimSpace(line); len(line) > 0 {
			r.logger.Errorf("%s: reader: unmarshal json error(%v): %s", r.collectorId, err, line)
		}
		return nil
	}
}

// 把文件seek到上次保存的位置
func (r *Reader) seekToSavedOffset() {
	if offset := r.offsetFile.read(); offset > 0 {
		if size := r.size(); size >= 0 && offset <= size {
			if _, err := r.file.Seek(offset, os.SEEK_SET); err == nil {
				r.logger.Infof("reader: seekToSavedOffset: %d, size: %d", offset, size)
			} else {
				r.logger.Errorf("%s: reader: seekToSavedOffset: %d, size: %d, error: %v",
					r.collectorId, offset, size, err,
				)
			}
		} else {
			r.logger.Errorf("%s: reader: seekToSavedOffset: offset %d exceeds size %d",
				r.collectorId, offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (r *Reader) seekFrontIfTruncated() {
	if offset := r.offset(); offset > 0 {
		if size := r.size(); size >= 0 && offset > size {
			if _, err := r.file.Seek(0, os.SEEK_SET); err == nil {
				r.logger.Infof("reader: seekFront(offset: %d, size: %d)", offset, size)
				r.reader.Reset(r.file)
			} else {
				r.logger.Errorf("%s: reader: seekFront(offset: %d, size: %d) error: %v",
					r.collectorId, offset, size, err,
				)
			}
		}
	}
}

func (r *Reader) offset() int64 {
	if offset, err := r.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset - int64(r.reader.Buffered()+len(r.buffered))
	} else {
		r.logger.Errorf("%s: reader: get offset error: %v", r.collectorId, err)
		return -1
	}
}

func (r *Reader) size() int64 {
	if fi, err := r.file.Stat(); err == nil {
		return fi.Size()
	} else {
		r.logger.Errorf("%s: reader: get size error: %v", r.collectorId, err)
		return -1
	}
}
