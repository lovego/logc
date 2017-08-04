package source

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/lovego/xiaomei/utils/fs"
)

// 把文件seek到上次保存的位置
func (s *Source) seekToSavedOffset() {
	if offset := s.readOffset(); offset > 0 {
		if size := s.size(); size >= 0 && offset <= size {
			if _, err := s.file.Seek(offset, os.SEEK_SET); err == nil {
				s.logger.Printf("seekToSavedOffset %s: offset(%d) size(%d)\n", s.path, offset, size)
			} else {
				s.logger.Printf("seekToSavedOffset %s error: %v\n", s.path, err)
			}
		} else {
			s.logger.Printf("seekToSavedOffset %s: offset(%d) exceeds size(%d)\n", s.path, offset, size)
		}
	}
}

// 如果文件被截短，把文件seek到开头
func (s *Source) seekFrontIfTruncated() {
	if offset := s.offset(); offset > 0 {
		if size := s.size(); size >= 0 && offset > size {
			if _, err := s.file.Seek(0, os.SEEK_SET); err == nil {
				s.logger.Printf("seekFront (size: %d, offset: %d)\n", size, offset)
				s.reader.Reset(s.file)
			} else {
				s.logger.Printf("seekFront (size: %d, offset: %d) error: %v\n", size, offset, err)
			}
		}
	}
}

func (s *Source) readOffset() int64 {
	if fs.NotExist(s.offsetPath) {
		return 0
	}
	content, err := ioutil.ReadFile(s.offsetPath)
	if err != nil {
		s.logger.Printf("read offset %s: %v\n", s.offsetPath, err)
		return 0
	}
	contentStr := strings.TrimSpace(string(content))
	if len(contentStr) == 0 {
		return 0
	}
	offset, err := strconv.ParseInt(contentStr, 10, 64)
	if err != nil {
		s.logger.Printf("parse offset %s(%s): %s\n", s.offsetPath, contentStr, err)
		return 0
	}
	return offset
}

func (s *Source) offset() int64 {
	if offset, err := s.file.Seek(0, os.SEEK_CUR); err == nil {
		return offset
	} else {
		s.logger.Printf("get offset error: %v\n", err)
		return -1
	}
}

func (s *Source) size() int64 {
	if fi, err := s.file.Stat(); err == nil {
		return fi.Size()
	} else {
		s.logger.Printf("get size error: %v\n", err)
		return -1
	}
}
