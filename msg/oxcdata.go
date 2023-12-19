package msg

import (
	"bytes"
	"io"
	"math"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-oxcdata/0c77892e-288e-435a-9c49-be1c20c7afdb
func GetDataValue(nametype string, data []byte) interface{} {
	switch nametype {
	case "0x0000", "PtypUnspecified":
		return PtypUnspecified(data)
	case "0x0001", "PtypNull":
		return PtypNull(data)
	case "0x0002", "PtypInteger16":
		return PtypInteger16(data)
	case "0x0003", "PtypInteger32":
		return PtypInteger32(data)
	case "0x0004", "PtypFloating32":
		return PtypFloating32(data)
	case "0x0005", "PtypFloating64":
		return PtypFloating64(data)
	case "0x0006", "PtypCurrency":
		return PtypCurrency(data)
	case "0x0007", "PtypFloatingTime":
		return PtypFloatingTime(data)
	case "0x000A", "PtypErrorCode":
		return PtypErrorCode(data)
	case "0x000B", "PtypBoolean":
		return PtypBoolean(data)
	case "0x000D", "PtypObject":
		return PtypObject(data)
	case "0x0014", "PtypInteger64":
		return PtypInteger64(data)
	case "0x001E", "PtypString8":
		return PtypString8(data)
	case "0x001F", "PtypString":
		return PtypString(data)
	case "0x0040", "PtypTime":
		return PtypTime(data)
	case "0x0048", "PtypGuid":
		return PtypGuid(data)
	case "0x00FB", "PtypServerId":
		return PtypServerId(data)
	case "0x00FD", "PtypRestriction":
		return PtypRestriction(data)
	case "0x00FE", "PtypRuleAction":
		return PtypRuleAction(data)
	case "0x0102", "PtypBinary":
		return PtypBinary(data)
	case "0x1002", "PtypMultipleInteger16":
		return PtypMultipleInteger16(data)
	case "0x1003", "PtypMultipleInteger32":
		return PtypMultipleInteger32(data)
	case "0x1004", "PtypMultipleFloating32":
		return PtypMultipleFloating32(data)
	case "0x1005", "PtypMultipleFloating64":
		return PtypMultipleFloating64(data)
	case "0x1006", "PtypMultipleCurrency":
		return PtypMultipleCurrency(data)
	case "0x1007", "PtypMultipleFloatingTime":
		return PtypMultipleFloatingTime(data)
	case "0x1014", "PtypMultipleInteger64":
		return PtypMultipleInteger64(data)
	case "0x101F", "PtypMultipleString":
		return PtypMultipleString(data)
	case "0x101E", "PtypMultipleString8":
		return PtypMultipleString8(data)
	case "0x1040", "PtypMultipleTime":
		return PtypMultipleTime(data)
	case "0x1048", "PtypMultipleGuid":
		return PtypMultipleGuid(data)
	case "0x1102", "PtypMultipleBinary":
		return PtypMultipleBinary(data)
	default:
		return data
	}
}

// this property type value matches any type
func PtypUnspecified(data []byte) []byte {
	return data
}

// this property is a placeholder
func PtypNull(data []byte) []byte {
	return nil
}

// 2 bytes; a 16-bit integer
func PtypInteger16(data []byte) uint16 {
	if len(data) >= 2 {
		return DecodeUint16(data[:2])
	}
	return 0
}

// 4 bytes; a 32-bit integer
func PtypInteger32(data []byte) uint32 {
	if len(data) >= 4 {
		return DecodeUint32(data[:4])
	}
	return 0
}

// 4 bytes; a 32-bit floating point number
func PtypFloating32(data []byte) float32 {
	if len(data) >= 4 {
		bits := DecodeUint32(data[:4])
		return math.Float32frombits(bits)
	}
	return 0
}

// 8 bytes; a 64-bit floating point number
func PtypFloating64(data []byte) float64 {
	if len(data) >= 8 {
		bits := DecodeUint64(data[:8])
		return math.Float64frombits(bits)
	}
	return 0
}

// 8 bytes;
// a 64-bit signed, scaled integer representation of a decimal currency value, with four places to the right of the decimal point
func PtypCurrency(data []byte) []byte {
	return data
}

// 8 bytes;
// a 64-bit floating point number
func PtypFloatingTime(data []byte) []byte {
	return data
}

// 4 bytes;
// A 32-bit integer encoding error information
func PtypErrorCode(data []byte) uint32 {
	if len(data) >= 4 {
		return DecodeUint32(data[:4])
	}
	return 0
}

// 1 byte; restricted to 1 or 0
func PtypBoolean(data []byte) bool {
	if len(data) > 0 && data[0] != 0 {
		return true
	}
	return false
}

// the property value is a Component Object Model (COM) object
func PtypObject(data []byte) []byte {
	if len(data) > 0 {
		return Trim(data)
	}
	return data
}

// 8 bytes; a 64-bit integer
func PtypInteger64(data []byte) uint64 {
	if len(data) >= 8 {
		return DecodeUint64(data[:8])
	}
	return 0
}

// variable size;
// a string of multibyte characters in externally specified encoding with terminating null character (single 0 byte).
func PtypString8(data []byte) string {
	if len(data) > 0 {

		reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder())

		utf8Data, err := io.ReadAll(reader)
		if err != nil {
			return string(Trim(data))
		}
		return string(Trim(utf8Data))
	}
	return string(data)
}

// variable size;
// a string of Unicode characters in UTF-16LE format encoding with terminating null character (0x0000).
func PtypString(data []byte) string {
	if len(data) > 0 {
		return UTF16ToUTF8(data)
	}
	return string(data)
}

// 8 bytes;
// a 64-bit integer representing the number of 100-nanosecond intervals since January 1, 1601
func PtypTime(data []byte) uint64 {
	if len(data) >= 8 {
		return DecodeUint64(data[:8])
	}
	return 0
}

// 16 bytes;
// a GUID with Data1, Data2, and Data3 fields in little-endian format
func PtypGuid(data []byte) []byte {
	return data
}

// variable size;
// a 16-bit COUNT field followed by a structure
func PtypServerId(data []byte) []byte {
	return data
}

// variable size;
// a byte array representing one or more Restriction structures
func PtypRestriction(data []byte) []byte {
	return data
}

// variable size;
// a 16-bit COUNT field followed by that many rule action structures
func PtypRuleAction(data []byte) []byte {
	return data
}

// variable size;
// a COUNT field followed by that many bytes
func PtypBinary(data []byte) []byte {
	if len(data) > 0 {
		return Trim(data)
	}
	return data
}

// variable size;
// a COUNT field followed by that many PtypInteger16 values
func PtypMultipleInteger16(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypInteger32 values
func PtypMultipleInteger32(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypFloating32 values
func PtypMultipleFloating32(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypFloating64 values
func PtypMultipleFloating64(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypCurrency values
func PtypMultipleCurrency(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypFloatingTime values
func PtypMultipleFloatingTime(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypInteger64 values
func PtypMultipleInteger64(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypString values
func PtypMultipleString(data []byte) interface{} {
	if len(data) > 0 {
		return string(Trim(data))
	}
	return data
}

// variable size;
// a COUNT field followed by that many PtypString8 values
func PtypMultipleString8(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypTime values
func PtypMultipleTime(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypGuid values
func PtypMultipleGuid(data []byte) interface{} {
	return data
}

// variable size;
// a COUNT field followed by that many PtypBinary values
func PtypMultipleBinary(data []byte) interface{} {
	return data
}
