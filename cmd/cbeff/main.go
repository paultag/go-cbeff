package main

import (
	"bytes"
	"encoding/asn1"
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

	rv := asn1.RawValue{}
	rest, err := asn1.Unmarshal(bytez, &rv)
	ohshit(err)

	if len(rest) != 0 {
		panic("Trailing garbage")
	}

	rvN := asn1.RawValue{}
	_, err = asn1.Unmarshal(rv.Bytes, &rvN)
	ohshit(err)

	c, err := cbeff.Parse(bytes.NewReader(rvN.Bytes))
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
