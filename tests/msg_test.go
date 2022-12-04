package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/mel2oo/mailfile/msg"
)

func TestParseMsg(t *testing.T) {
	msg, err := msg.New("testdata/complete.msg")
	if err != nil {
		t.Fail()
		return
	}

	for _, attach := range msg.Attachments {
		data, err := io.ReadAll(attach.Data)
		if err != nil {
			continue
		}

		fmt.Printf("name:%s md5:%s\n",
			attach.Filename, GetMD5(data))
	}

	msg.Message.Output()
}

func GetMD5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
