package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Source struct {
	path       string
	offsetPath string
	logger     *log.Logger
	reader     *bufio.Reader
	file       *os.File
}

func New(path, offsetPath string, logger *log.Logger) *Source {
	return &Source{
		path: path, offsetPath: offsetPath, logger: logger,
	}
}

func (s *Source) Read() (rows []map[string]interface{}, drain bool) {
	if !s.setupFileAndReader() {
		return nil, true
	}
	for size := 0; size < 2*1024*1024; {
		line, err := s.reader.ReadBytes('\n')
		if row := s.parseRow(line); row != nil {
			rows = append(rows, row)
			size += len(line)
		}
		if err != nil {
			if err == io.EOF {
				if size == 0 && len(line) == 0 {
					s.seekFrontIfTruncated()
				}
			} else {
				s.logger.Println(`read error: ` + err.Error())
			}
			return rows, true
		}
	}
	return rows, false
}

func (s *Source) SaveOffset() string {
	var offsetStr string
	if offset := s.offset(); offset > 0 {
		offsetStr = strconv.FormatInt(offset, 10)
	} else {
		return ``
	}

	if err := ioutil.WriteFile(s.offsetPath, []byte(offsetStr), 0666); err != nil {
		s.logger.Printf("write offset error: %v\n", err.Error())
	}
	return offsetStr
}

func (s *Source) Opened() bool {
	return s.file != nil
}

func (s *Source) Reopen() {
	if s.file != nil {
		s.file.Close()
		s.file = nil
		s.logger.Printf("reopen %s", s.path)
	}
	if err := os.Remove(s.offsetPath); err != nil && !os.IsNotExist(err) {
		s.logger.Printf("remove offsetFile %s error: %v\n", s.offsetPath, err.Error())
	}
}

func (s *Source) parseRow(line []byte) map[string]interface{} {
	if len(line) == 0 {
		return nil
	}
	var row map[string]interface{}
	if err := json.Unmarshal(line, &row); err == nil {
		return row
	} else {
		if line = bytes.TrimSpace(line); len(line) > 0 {
			s.logger.Printf("json error(%v): %s\n", err, line)
		}
		return nil
	}
}

func (s *Source) setupFileAndReader() bool {
	if s.file != nil {
		return true
	}
	if file, err := os.Open(s.path); err != nil {
		s.logger.Printf("open %s: %v\n", s.path, err)
		return false
	} else {
		s.file = file
		s.seekToSavedOffset()
		s.reader = bufio.NewReader(s.file)
		return true
	}
}
