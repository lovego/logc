package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open(`t.txt`)
	if err != nil {
		log.Panic(err)
	}
	for {
		printOffset(f)
		readData(f)
		printOffset(f)

		waitInput()
	}
}

func printOffset(f *os.File) {
	if offset, err := f.Seek(0, os.SEEK_CUR); err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("offset: %d\n", offset)
	}
}

func readData(f *os.File) {
	buf := make([]byte, 7)
	if n, err := f.Read(buf); err != nil && err != io.EOF {
		log.Panic(err)
	} else {
		fmt.Printf("data: %d %#v\n", n, string(buf[:n]))
	}
}

var stdinReader = bufio.NewReader(os.Stdin)

func waitInput() {
	if _, err := stdinReader.ReadString('\n'); err != nil {
		log.Panic(err)
	}
}
