package cbeff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	//"os"
)

type Time [8]byte

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
	BiometricType         [3]byte
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

type FacialHeader struct {
	FormatID     [4]byte
	VersionID    [4]byte
	RecordLength uint32
	NumberFaces  uint16
}

func (fh FacialHeader) Validate() error {
	if bytes.Compare(fh.FormatID[:], []byte{'F', 'A', 'C', 0x00}) != 0 {
		return fmt.Errorf("cbeff: FacialHeader.FormatID isn't FAC\\0")
	}

	if bytes.Compare(fh.VersionID[:], []byte{'0', '1', '0', 0x00}) != 0 {
		return fmt.Errorf("cbeff: FacialHeader.VersionID isn't 010\\0")
	}

	if fh.NumberFaces != 1 {
		return fmt.Errorf(
			"cbeff: FacialHeader.NumberFaces isn't 1, and I got confused",
		)
	}

	return nil
}

type FacialInformation struct {
	Length                  uint32
	NumberOfPoints          uint16
	BiographicalInformation [6]byte
	Expression              [2]byte
	Pose                    [3]byte
	PoseUncertainty         [3]byte
}

func (fi FacialInformation) Validate() error {
	// Not currently checking for anything.
	if fi.NumberOfPoints != 0x00 {
		return fmt.Errorf("cbeff: FacialInformation.NumberOfPoints isn't 0")
	}
	return nil
}

type ImageInformation struct {
	Type       uint8
	DataType   uint8
	Width      uint16
	Height     uint16
	ColorSpace uint8
	SourceType uint8
	DeviceType uint16
	Quality    uint16
}

func (ii ImageInformation) Validate() error {
	if ii.Type != 0x01 {
		return fmt.Errorf("cbeff: ImageInformation.Type isn't 1")
	}
	return nil
}

func Parse(in io.Reader) (*Header, error) {
	h := Header{}
	if err := binary.Read(in, binary.BigEndian, &h); err != nil {
		return nil, err
	}
	if err := h.Validate(); err != nil {
		return nil, err
	}

	fh := FacialHeader{}
	if err := binary.Read(in, binary.BigEndian, &fh); err != nil {
		return nil, err
	}
	if err := fh.Validate(); err != nil {
		return nil, err
	}

	fi := FacialInformation{}
	if err := binary.Read(in, binary.BigEndian, &fi); err != nil {
		return nil, err
	}
	if err := fi.Validate(); err != nil {
		return nil, err
	}

	ii := ImageInformation{}
	if err := binary.Read(in, binary.BigEndian, &ii); err != nil {
		return nil, err
	}
	if err := ii.Validate(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(io.LimitReader(in, int64(fi.Length)))
	if err != nil {
		return nil, err
	}

	// fd, err := os.Create("output.bin")
	// if err != nil {
	// 	return nil, err
	// }
	// defer fd.Close()
	// fd.Write(data)

	_ = data
	_ = fmt.Printf

	return nil, nil
}
