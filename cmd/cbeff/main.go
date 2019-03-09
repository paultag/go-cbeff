package main

import (
	"fmt"
	"os"

	"pault.ag/go/cbeff"
)

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fd, err := os.Open(os.Args[1])
	ohshit(err)

	data, err := cbeff.Parse(fd)
	ohshit(err)

	_ = data
	_ = fmt.Printf
}
