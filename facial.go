package cbeff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

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
