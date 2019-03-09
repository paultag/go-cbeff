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

	h, err := cbeff.Parse(bytes.NewReader(rvN.Bytes))
	ohshit(err)

	_ = h
	_ = fmt.Printf
}
