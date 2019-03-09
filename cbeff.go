package cbeff

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Options uint8

func (o Options) BitSet(bit uint8) bool       { return (uint8(o)&bit == bit) }
func (o Options) OptionalFieldsPresent() bool { return o.BitSet(0x08) }
func (o Options) SignatureOrMAC() bool        { return o.BitSet(0x04) }
func (o Options) Privacy() bool               { return o.BitSet(0x02) }
func (o Options) Integrity() bool             { return o.BitSet(0x01) }

type OptionalField uint16

func (of OptionalField) String() string {
	e, ok := FieldValueName[uint16(of)]
	if !ok {
		return "unknown"
	}
	return e
}

func (of OptionalField) Fields() []OptionalField {
	ret := []OptionalField{}
	for _, field := range AllOptionalFields {
		fi := uint16(field)
		if (uint16(of) & fi) == fi {
			ret = append(ret, field)
		}
	}
	return ret
}

const (
	OptionalFieldSubheaderCount         OptionalField = 0x0001
	OptionalFieldPID                    OptionalField = 0x0002
	OptionalFieldPatronFormatIdentifier OptionalField = 0x0004
	OptionalFieldIndex                  OptionalField = 0x0008

	OptionalFieldBiometricCreationDate OptionalField = 0x0010
	OptionalFieldValidityPeriod        OptionalField = 0x0020
	OptionalFieldBiometricType         OptionalField = 0x0040
	OptionalFieldBiometricSubtype      OptionalField = 0x0080

	OptionalFieldCBEFFHeaderVersion  OptionalField = 0x100
	OptionalFieldPatronHeaderVersion OptionalField = 0x200
	OptionalFieldBiometricPurpose    OptionalField = 0x400
	OptionalFieldBiometricDataType   OptionalField = 0x800

	OptionalFieldBiometricDataQuality OptionalField = 0x1000
	OptionalFieldCreator              OptionalField = 0x2000
	OptionalFieldChallengeResponse    OptionalField = 0x4000
	OptionalFieldPayload              OptionalField = 0x8000
)

func Reverse(e map[string]uint16) map[uint16]string {
	ret := map[uint16]string{}
	for k, v := range e {
		ret[v] = k
	}
	return ret
}

var (
	AllOptionalFields = []OptionalField{
		OptionalFieldSubheaderCount, OptionalFieldPID,
		OptionalFieldPatronFormatIdentifier, OptionalFieldIndex,
		OptionalFieldBiometricCreationDate, OptionalFieldValidityPeriod,
		OptionalFieldBiometricType, OptionalFieldBiometricSubtype,
		OptionalFieldCBEFFHeaderVersion, OptionalFieldPatronHeaderVersion,
		OptionalFieldBiometricPurpose, OptionalFieldBiometricDataType,
		OptionalFieldBiometricDataQuality, OptionalFieldCreator,
		OptionalFieldChallengeResponse, OptionalFieldPayload,
	}

	FieldNameValue = map[string]uint16{
		"SubheaderCount":         0x0001,
		"PID":                    0x0002,
		"PatronFormatIdentifier": 0x0004,
		"Index":                  0x0008,
		"BiometricCreationDate":  0x0010,
		"ValidityPeriod":         0x0020,
		"BiometricType":          0x0040,
		"BiometricSubtype":       0x0080,
		"CBEFFHeaderVersion":     0x100,
		"PatronHeaderVersion":    0x200,
		"BiometricPurpose":       0x400,
		"BiometricDataType":      0x800,
		"BiometricDataQuality":   0x1000,
		"Creator":                0x2000,
		"ChallengeResponse":      0x4000,
		"Payload":                0x8000,
	}
	FieldValueName map[uint16]string = Reverse(FieldNameValue)
)

type Header struct {
	BDBFormatOwner [2]byte
	BDBFormatType  [2]byte
	Options        Options
}

func Parse(in io.Reader) (*Header, error) {
	data := Header{}
	if err := binary.Read(in, binary.LittleEndian, &data); err != nil {
		return nil, err
	}

	if data.Options.OptionalFieldsPresent() {
		of := OptionalField(0)
		if err := binary.Read(in, binary.LittleEndian, &of); err != nil {
			return nil, err
		}
		fmt.Printf("%s\n", of.Fields())
	}

	return nil, nil
}

// typedef struct bioapi_bir {
// BioAPI_BIR_HEADER Header;
// BioAPI_BIR_BIOMETRIC_DATA_PTR BiometricData; /* length indicated in header */
// BioAPI_DATA_PTR Signature; /* NULL if no signature; length is inherent in this type */
// } BioAPI_BIR, *BioAPI_BIR_PTR;
//
// typedef struct bioapi_bir_header {
// uint32 Length; /* Length of Header + Opaque Data */
// BioAPI_BIR_VERSION HeaderVersion;
// BioAPI_BIR_DATA_TYPE Type;
// BioAPI_BIR_BIOMETRIC_DATA_FORMAT Format;
// BioAPI_QUALITY Quality;
// BioAPI_BIR_PURPOSE PurposeMask;
// BioAPI_BIR_AUTH_FACTORS FactorsMask;
// } BioAPI_BIR_HEADER, *BioAPI_BIR_HEADER_PTR;
//
// typedef struct bioapi_bir_biometric_data_format {
// uint16 FormatOwner;
// uint16 FormatID;
// } BioAPI_BIR_BIOMETRIC_DATA_FORMAT, *BioAPI_BIR_BIOMETRIC_DATA_FORMAT_PTR;
//
// typedef uint8 BioAPI_BIR_BIOMETRIC_DATA;
