package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open(`t.txt`)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	for {
		printOffset(f)
		readData(reader)
		printOffset(f)

		waitInput()
	}
}

func printOffset(f *os.File) {
	if offset, err := f.Seek(0, os.SEEK_CUR); err != nil {
		panic(err)
	} else {
		fmt.Printf("offset: %d\n", offset)
	}
}

func readData(reader *bufio.Reader) {
	if line, err := reader.ReadString('\n'); err != nil && err != io.EOF {
		panic(err)
	} else {
		fmt.Printf("data: %d %#v\n", len(line), line)
	}
}

var stdinReader = bufio.NewReader(os.Stdin)

func waitInput() {
	if _, err := stdinReader.ReadString('\n'); err != nil {
		panic(err)
	}
}
