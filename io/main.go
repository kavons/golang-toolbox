package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

func main() {
    file, err := os.Open("./input")
    if err != nil {
    	log.Fatal(err)
	}

	buf := make([]byte, 10)
	n, err := io.ReadFull(file, buf)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(n, string(buf[:]))
	}

	m := int64(binary.LittleEndian.Uint64(buf))
	log.Println(m)
}