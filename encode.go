package msgfile

import (
	"bytes"
	"encoding/binary"
	"strings"

	"golang.org/x/text/encoding/unicode"
)

func UTF16ToUTF8(b []byte) string {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	data, err := decoder.Bytes(b)
	if err != nil {
		return string(b)
	}
	return string(data)
}

func Trim(data []byte) []byte {
	return []byte(strings.ReplaceAll(string(data), "\x00", ""))
}

//FromUnicode read unicode and convert to byte array
func FromUnicode(uni []byte) string {
	st := ""
	for _, k := range uni {
		if k != 0x00 {
			st += string(k)
		}
	}
	return st
}

//DecodeInt64 decode 8 byte value into int64
func DecodeInt64(num []byte) int64 {
	var number int64
	bf := bytes.NewReader(num)
	binary.Read(bf, binary.BigEndian, &number)
	return number
}

//DecodeUint64 decode 4 byte value into uint32
func DecodeUint64(num []byte) uint64 {
	var number uint64
	bf := bytes.NewReader(num)
	binary.Read(bf, binary.LittleEndian, &number)
	return number
}

//DecodeUint32 decode 4 byte value into uint32
func DecodeUint32(num []byte) uint32 {
	var number uint32
	bf := bytes.NewReader(num)
	binary.Read(bf, binary.LittleEndian, &number)
	return number
}

//DecodeUint16 decode 2 byte value into uint16
func DecodeUint16(num []byte) uint16 {
	var number uint16
	bf := bytes.NewReader(num)
	binary.Read(bf, binary.LittleEndian, &number)
	return number
}

//DecodeUint8 decode 1 byte value into uint8
func DecodeUint8(num []byte) uint8 {
	var number uint8
	bf := bytes.NewReader(num)
	binary.Read(bf, binary.LittleEndian, &number)
	return number
}

//ReadUint32 read 4 bytes and return as uint32
func ReadUint32(pos int, buff []byte) (uint32, int) {
	return DecodeUint32(buff[pos : pos+4]), pos + 4
}

//ReadUint16 read 2 bytes and return as uint16
func ReadUint16(pos int, buff []byte) (uint16, int) {
	return DecodeUint16(buff[pos : pos+2]), pos + 2
}

//ReadUint8 read 1 byte and return as uint8
func ReadUint8(pos int, buff []byte) (uint8, int) {
	return DecodeUint8(buff[pos : pos+2]), pos + 2
}

//ReadBytes read and return count number o bytes
func ReadBytes(pos, count int, buff []byte) ([]byte, int) {
	return buff[pos : pos+count], pos + count
}

//ReadByte read and return a single byte
func ReadByte(pos int, buff []byte) (byte, int) {
	return buff[pos : pos+1][0], pos + 1
}

//ReadUnicodeString read and return a unicode string
func ReadUnicodeString(pos int, buff []byte) ([]byte, int) {
	//stupid hack as using bufio and ReadString(byte) would terminate too early
	//would terminate on 0x00 instead of 0x0000
	index := bytes.Index(buff[pos:], []byte{0x00, 0x00})
	if index == -1 {
		return nil, 0
	}
	str := buff[pos : pos+index]
	return []byte(str), pos + index + 2
}

//ReadUTF16BE reads the unicode string that the outlook rule file uses
//this basically means there is a length byte that we need to skip over
func ReadUTF16BE(pos int, buff []byte) ([]byte, int) {

	lenb := (buff[pos : pos+1])
	k := int(lenb[0])
	pos += 1 //length byte but we don't really need this
	var str []byte
	if k == 0 {
		str, pos = ReadUnicodeString(pos, buff)
	} else {
		str, pos = ReadBytes(pos, k*2, buff) //
		//pos += 2
	}

	return str, pos
}

//ReadASCIIString returns a string as ascii
func ReadASCIIString(pos int, buff []byte) ([]byte, int) {
	bf := bytes.NewBuffer(buff[pos:])
	str, _ := bf.ReadString(0x00)
	return []byte(str), pos + len(str)
}

//ReadTypedString reads a string as either Unicode or ASCII depending on type value
func ReadTypedString(pos int, buff []byte) ([]byte, int) {
	var t = buff[pos]
	if t == 0 { //no string
		return []byte{}, pos + 1
	}
	if t == 1 {
		return []byte{}, pos + 1
	}
	if t == 3 {
		str, p := ReadASCIIString(pos+1, buff)
		return str, p
	}
	if t == 4 {
		str, p := ReadUnicodeString(pos+1, buff)
		return str, p
	}
	str, _ := ReadBytes(pos+1, 4, buff)
	return str, pos + len(str)
}
