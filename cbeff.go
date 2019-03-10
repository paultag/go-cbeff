package cbeff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
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

type Facial struct {
	Header FacialHeader
	Reader io.Reader
	Images []Image
}

func (c Facial) nextImage() (*Image, error) {
	fi := FacialInformation{}
	if err := binary.Read(c.Reader, binary.BigEndian, &fi); err != nil {
		return nil, err
	}
	if err := fi.Validate(); err != nil {
		return nil, err
	}

	ii := ImageInformation{}
	if err := binary.Read(c.Reader, binary.BigEndian, &ii); err != nil {
		return nil, err
	}
	if err := ii.Validate(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(io.LimitReader(c.Reader, int64(fi.Length)))
	if err != nil {
		return nil, err
	}

	return &Image{
		FacialInformation: fi,
		ImageInformation:  ii,
		Data:              data,
	}, nil
}

func (c CBEFF) Facial() (*Facial, error) {
	if !c.Header.BiometricType.Equal(BiometricTypeFacial) {
		return nil, fmt.Errorf("cbeff: Header.BiometricType isn't Facial")
	}

	fh := FacialHeader{}
	if err := binary.Read(c.Reader, binary.BigEndian, &fh); err != nil {
		return nil, err
	}
	if err := fh.Validate(); err != nil {
		return nil, err
	}

	if fh.RecordLength != (c.Header.BDBLength) {
		return nil, fmt.Errorf(
			"cbeff: FacialHeader length disagrees with CBEFF length",
		)
	}

	f := Facial{
		Header: fh,
		Images: []Image{},
		Reader: io.LimitReader(c.Reader, int64(fh.RecordLength)),
	}

	for {
		image, err := f.nextImage()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		f.Images = append(f.Images, *image)
	}

	return &f, nil
}

func (h Header) Validate() error {
	if h.PatronHeaderVersion != 0x03 {
		return fmt.Errorf("cbeff: Header.PatronHeaderVersion isn't 3")
	}
	return nil
}

type Image struct {
	FacialInformation FacialInformation
	ImageInformation  ImageInformation
	Data              []byte
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
