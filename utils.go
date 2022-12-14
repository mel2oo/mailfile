package mailfile

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/djimenez/iconv-go"
)

//解析 带有 =?utf-8?B?bWxlbW9z?= 格式字符串
func ParseContext(data string) string {
	bindex := strings.Index(data, "=?")
	if bindex == -1 {
		return data
	}

	eindex := strings.Index(data[bindex:], "?=")
	if eindex == -1 {
		return data
	}

	if bindex+eindex+2 > len(data) {
		return data
	}

	prefix := data[:bindex]
	decode := ParseTitle(data[bindex : bindex+eindex+2])
	suffix := ParseContext(data[bindex+eindex+2:])
	return prefix + decode + suffix
}

func ParseTitle(subject string) string {
	if !strings.HasPrefix(subject, "=?") {
		return subject
	}

	title := subject[2 : len(subject)-2]
	lists := strings.SplitN(title, "?", 3)
	if len(lists) != 3 {
		return subject
	}

	var charset string
	if lists[1] == "b" || lists[1] == "B" {
		charset = "base64"
	} else if lists[1] == "q" || lists[1] == "Q" {
		charset = "quoted-printable"
	} else {
		return subject
	}

	texts, err := DecodeString(lists[2], charset)
	if err != nil {
		return subject
	}

	retstr, err := ConvertData(texts, lists[0])
	if err == nil {
		return retstr
	}

	if strings.ToLower(lists[0]) == "gb2312" {
		data, err := iconv.ConvertString(string(texts), "gb2312", "utf-8")
		if err != nil {
			return subject
		}
		return data
	}
	return subject
}

func DecodeString(str string, etype string) ([]byte, error) {
	switch strings.ToLower(etype) {
	case "":
		return []byte(str), nil
	case "base64":
		return base64.StdEncoding.DecodeString(str)
	case "quoted-printable":
		reader := strings.NewReader(str)
		decode := quotedprintable.NewReader(reader)
		return ioutil.ReadAll(decode)
	default:
		return nil, fmt.Errorf("unkown encode format %s", etype)
	}
}

func ConvertData(data []byte, charset string) (string, error) {
	charset = strings.ToLower(charset)
	switch strings.ToLower(charset) {
	case "utf-8", "utf8":
		return string(data), nil
	default:
		data, err := iconv.ConvertString(string(data), charset, "utf-8")
		if err == nil {
			return data, err
		}

		decoder1 := mahonia.NewDecoder(charset)
		if decoder1 == nil {
			return string(data), fmt.Errorf("not found decode %s", charset)
		}

		decoder2 := mahonia.NewDecoder("utf-8")
		if decoder2 == nil {
			return string(data), fmt.Errorf("not found decode %s", charset)
		}

		_, cdata, err := decoder2.Translate([]byte(
			decoder1.ConvertString(string(data))), true)
		return string(cdata), err
	}
}
