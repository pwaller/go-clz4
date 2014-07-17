package main

import (
	"flag"
	// "fmt"
	"io"
	"log"
	"os"

	"github.com/pwaller/go-clz4/stream"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalln("usage: test [filename.lz4]")
	}
	name := flag.Arg(0)

	fd, err := os.Open(name)
	if err != nil {
		log.Fatalln("Failed to open file", name, err)
	}
	defer fd.Close()

	s, err := stream.NewReader(fd)
	if err != nil {
		log.Fatalln("Error parsing stream header", err)
	}

	n, err := io.Copy(os.Stdout, s)
	log.Println(n, err)

	// p := make([]byte, 1024)

	// n, err := s.Read(p)
	// log.Println(n, err)
	// fmt.Printf("%s", p)
	_ = s
}
