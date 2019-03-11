package main

import (
	"fmt"
	"io/ioutil"
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

	bytez, err := ioutil.ReadAll(fd)
	ohshit(err)

	c, err := cbeff.ParsePIV(bytez)
	ohshit(err)

	f, err := c.Facial()
	ohshit(err)

	for _, image := range f.Images {
		fd, err := os.Create("test.j2")
		ohshit(err)
		defer fd.Close()
		fd.Write(image.Data)
		fd.Close()
	}

	_ = fmt.Printf
}
