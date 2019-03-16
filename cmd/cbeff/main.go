package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"pault.ag/go/cbeff"
	"pault.ag/go/cbeff/jpeg2000"
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
		fmt.Printf("%s\n", image.FacialInformation.BiographicalInformation)

		img, err := jpeg2000.Parse(image.Data)
		ohshit(err)

		fd, err := os.Create("test.png")
		ohshit(err)
		defer fd.Close()
		if err := png.Encode(fd, img); err != nil {
			fd.Close()
			ohshit(err)
		}
	}

	_ = fmt.Printf
}
