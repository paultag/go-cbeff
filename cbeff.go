package cbeff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type CBEFF struct {
	Header Header
	Reader io.Reader
}

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

type Time [8]byte
type BiometricType [3]byte

func (b BiometricType) Equal(o BiometricType) bool {
	return bytes.Compare(b[:], o[:]) == 0
}

var (
	BiometricTypeFingerprint = BiometricType{0x00, 0x00, 0x08}
	BiometricTypeFacial      = BiometricType{0x00, 0x00, 0x02}
)

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

func (h Header) Validate() error {
	if h.PatronHeaderVersion != 0x03 {
		return fmt.Errorf("cbeff: Header.PatronHeaderVersion isn't 3")
	}
	return nil
}
