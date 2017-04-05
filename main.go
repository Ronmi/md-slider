package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var fn string
	flag.StringVar(&fn, "i", "", "input filename, markdown format")
	flag.Parse()

	ret, err := conv(fn)
	if err != nil {
		log.Fatalf("Cannot convert %s: %s", fn, err)
	}

	fmt.Println(ret)
}
