package mailfile

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/djimenez/iconv-go"
)

// 解析 带有 =?utf-8?B?bWxlbW9z?= 格式字符串
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
	return strings.TrimSpace(prefix) + strings.TrimSpace(decode) + strings.TrimSpace(suffix)
}

func ParseFrom(from string) ([]*mail.Address, error) {
	if !strings.HasPrefix(from, "=?") {
		return mail.ParseAddressList(from)
	}
	idx := strings.LastIndex(from, "?=")
	if idx == -1 {
		return mail.ParseAddressList(from)
	}

	name := ParseTitle(from[:idx+2])
	from = name + from[idx+2:]
	return mail.ParseAddressList(from)
}

func ParseTitle(subject string) string {
	if strings.Count(subject, "?=") > 1 && len(strings.Split(subject, " ")) > 1 { // 有多个 =?[]?[]?[]?= 格式字符串
		var r strings.Builder
		subjects := strings.Split(subject, " ")
		for _, s := range subjects {
			s = ParseTitle(strings.TrimSpace(s))
			r.WriteString(s)
		}

		return r.String()
	}

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
		return io.ReadAll(decode)
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

func ParsePasswd(html, text []byte) []string {
	pwds := make(map[string]bool)

	if len(html) > 0 {
		ExtractPwd(TrimHTML(string(html)), &pwds)
	}

	if len(text) > 0 {
		ExtractPwd(string(text), &pwds)
	}

	res := make([]string, 0)
	for pwd := range pwds {
		pwd = strings.TrimSpace(pwd)
		if len(pwd) > 3 {
			res = append(res, pwd)
		}
	}
	return res
}

var (
	expHtml = regexp.MustCompile(`<[\s\S]*?>`)
	// ExpUnicode    = regexp.MustCompile(`&#\d+;`)
	// ExpUnicodeHex = regexp.MustCompile(`&#x\s+;`)

	expPasswd     = regexp.MustCompile(`[0-9a-zA-Z~$&+,:;=?@#|'<>.-^*()%!][0-9a-zA-Z~$&+,:;=?@#|'<>.-^*()%! ]{2,20}`)
	expPasswdUTF8 = regexp.MustCompile("[\u4e00-\u9fa50-9a-zA-Z~$&+,:;=?@#|'<>.-^*()%!][\u4e00-\u9fa50-9a-zA-Z~$&+,:;=?@#|'<>.-^*()%! ]{2,20}")

	keys = []string{
		"password",
		"passwd",
		"密码",
		"密码是",
		"密码为",
		"秘密",
		"パスワード",
		"Пароль ",
		"Pasvorto",
		"Mot de passe",
		"Passwort",
		"Contraseña",
		"कूटसङ्केतः",
		"암호",
	}
)

func TrimHTML(data string) string {

	txt := expHtml.ReplaceAllString(data, "")

	// unicode := expUnicode.FindAllString(txt, -1)
	// unicodeHex := expUnicodeHex.FindAllString(txt, -1)
	// filter := make(map[string]bool)

	// for _, code := range unicode {
	// 	if _, has := filter[code]; has {
	// 		continue
	// 	}
	// 	filter[code] = true
	// 	val, err := strconv.ParseInt(code[2:len(code)-1], 10, 0)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	txt = strings.ReplaceAll(txt, code, string(rune(val)))
	// }
	return txt
}

// rune占用大小
// func GetRuneLenth(org rune) int {
// 	if org < 128 {
// 		return 1
// 	} else if org < 2048 {
// 		return 2
// 	} else if org < 65536 {
// 		return 3
// 	} else {
// 		return 4
// 	}
// }

// 获取关键词最后一个字符串所在位置 当前关键词没有那种papasswd类型 可以不考虑kmp算法 之后存在可以考虑加上
// func FindStrLastIndex(source, key string) int {

// 	key_index := 0
// 	keys := []rune(key)
// 	for k, v := range []rune(source) {
// 		if v == keys[key_index] {
// 			key_index++
// 			if key_index == len(keys) {
// 				return k + GetRuneLenth(v)
// 			}
// 		} else if v == ' ' || v == '\t' || v == '\n' {
// 			continue
// 		} else {
// 			key_index = 0
// 		}
// 	}
// 	return -1
// }

// 字符串匹配
func SundaySearch(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}

	if len(haystack) == 0 && len(needle) != 0 {
		return -1
	}

	sIndex := 0
	pIndex := 0
	space := 0

	keys := make([]int, 256) //记录每个字符出现的最右侧位置

	for i := 0; i < 256; i++ {
		keys[i] = -1
	}

	for i := 0; i < len(needle); i++ {
		keys[needle[i]] = i
	}

	for sIndex < len(haystack) {
		if haystack[sIndex] == needle[pIndex] {
			sIndex++
			pIndex++
			if pIndex == len(needle) {
				return sIndex
			}
		} else if haystack[sIndex] == ' ' || haystack[sIndex] == '\t' || haystack[sIndex] == '\n' {
			sIndex++
			continue
		} else {
			pIndex = 0
			if space+len(needle) < len(haystack) {
				space += len(needle) - keys[haystack[space+len(needle)]] //寻找下个模式串匹配的初始位置
			} else { //位移与模式串的长度和超出原字符串的长度
				return -1
			}
			sIndex = space
		}
	}

	return -1
}
func ExtractPwd(data string, filter *map[string]bool) {
	var (
		tdata  string
		lowstr = strings.ToLower(data)
	)

	for _, key := range keys {
		lowkey := strings.ToLower(key)
		index := 0
		tdata = lowstr[index:]
		for find_index := SundaySearch(tdata, lowkey); find_index != -1; find_index = SundaySearch(tdata, lowkey) {
			index = index + find_index
			if index > len(lowstr) {
				break
			}
			if strings.HasPrefix(data[index:], "路径") {
				tdata = lowstr[index:]
				continue
			}
			pw := expPasswd.FindStringIndex(data[index:])
			if len(pw) > 1 {
				if pw[0] > 20 {
					tdata = lowstr[index:]
					continue
				}
				pwd := data[index:][pw[0]:pw[1]]
				pwd = strings.TrimPrefix(pwd, ":")
				(*filter)[strings.TrimSpace(pwd)] = true
			} else {
				pw = expPasswdUTF8.FindStringIndex(data[index:])
				if len(pw) > 1 {
					pwd := data[index:][pw[0]:pw[1]]
					pwd = strings.TrimPrefix(pwd, ":")
					(*filter)[strings.TrimSpace(pwd)] = true
				}
			}

			tdata = lowstr[index:]
		}
	}
}
