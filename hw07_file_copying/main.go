package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		log.Fatal("Usage: -from <source> -to <destination> [-offset N] [-limit N]")
	}

	if offset < 0 || limit < 0 {
		log.Fatal("offset and limit must be non-negative")
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Copy succeeded.")
}
