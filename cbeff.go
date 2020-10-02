// {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2019
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package cbeff // import "pault.ag/go/cbeff"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"pault.ag/go/fasc"
)

// CBEFF is an ecapsulation of a CBEFF serialized file. This will allow you
// to dispatch based on the Header.BiometricType, and extract the information
// depending on what type of CBEFF entry it is.
type CBEFF struct {
	// CBEFF file metadata, such as what kind of file this is, sizes of the
	// various payloads, and information such as creator, who the metric is of,
	// and validity time.
	Header Header

	// io.Reader over this CBEFF entry. This reader must be fully read before
	// moving onto any other files on the same Reader.
	Reader io.Reader
}

// Close the file.
func (c CBEFF) Close() error {
	_, err := io.Copy(ioutil.Discard, c.Reader)
	return err
}

// Parse will read the CBEFF data from the io.Reader, parse the CBEFF header,
// and construct he CBEFF file encapsulation.
func Parse(in io.Reader) (*CBEFF, error) {
	ret := CBEFF{}

	h := Header{}
	if err := binary.Read(in, binary.BigEndian, &h); err != nil {
		return nil, err
	}
	if err := h.Validate(); err != nil {
		return nil, err
	}

	ret.Header = h
	totalLength := int64(h.BDBLength) + int64(h.SBLength)
	ret.Reader = io.LimitReader(in, totalLength)

	return &ret, nil
}

// Time represents CBEFF Time as an 8 octet array, in the format of
// Y Y M D h m s Z, where Z is a literal ASCII 'Z', and the other values
// being the uint8 value for that position.
type Time [8]byte

// Time will return CBEFF Time into a Golang time.Time.
func (t Time) Time() (time.Time, error) {
	if t[7] != 'Z' {
		return time.Time{}, fmt.Errorf("cbeff: Time doesn't end with Z")
	}
	year := (int(t[0]) * 100) + int(t[1])
	month := time.Month(t[2])
	day := int(t[3])
	hour := int(t[4])
	minute := int(t[5])
	second := int(t[6])

	return time.Date(year, month, day, hour, minute, second, 0, time.UTC), nil
}

// BiometricType indicates the type of biometric stored in the CBEFF, such as
// Face photos, or Fingerprints.
type BiometricType [3]byte

// Equal will check to see if the two BiometricTypes are the same.
func (b BiometricType) Equal(o BiometricType) bool {
	return bytes.Compare(b[:], o[:]) == 0
}

var (
	// BiometricTypeFingerprint indicates the CBEFF file contains fingerprint
	// information. This may either be an enrollment or minutiae.
	BiometricTypeFingerprint = BiometricType{0x00, 0x00, 0x08}

	// BiometricTypeFacial indicates the CBEFF file contains the facial photos
	// to be used for visual confirmation of the individual.
	BiometricTypeFacial = BiometricType{0x00, 0x00, 0x02}
)

// Header contains information regarding the CBEFF data contained within
// the data stream.
type Header struct {
	PatronHeaderVersion   uint8
	SBHSecurityOptions    uint8
	BDBLength             uint32
	SBLength              uint16
	BDBFormatOwner        uint16
	BDBFormatType         uint16
	BiometricCreationDate Time
	ValidityNotBefore     Time
	ValidityNotAfter      Time
	BiometricType         BiometricType
	BiometricDataType     uint8
	BiometricDataQuality  uint8
	Creator               [18]byte
	FASC                  [25]byte
	Reserved              [4]byte
}

// ParseFASC will read the FASC bytes, and return a parsed pault.ag/go/fasc.Fasc
// struct.
func (h Header) ParseFASC() (*fasc.FASC, error) {
	return fasc.Parse(h.FASC[:])
}

// Validate will ensure that the header is understood by this library.
func (h Header) Validate() error {
	if h.PatronHeaderVersion != 0x03 {
		return fmt.Errorf("cbeff: Header.PatronHeaderVersion isn't 3")
	}
	return nil
}

// vim: foldmethod=marker
