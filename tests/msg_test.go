package test

import (
	"fmt"
	"testing"

	"github.com/mel2oo/mailfile/eml"
	"github.com/mel2oo/mailfile/msg"
)

func TestParseMsg(t *testing.T) {
	msg, err := msg.New("testdata/complete.msg")
	if err != nil {
		t.Fail()
		return
	}

	out := msg.Format()
	if len(out.Attachments) == 0 ||
		len(out.Embeddeds) == 0 ||
		len(out.SubMessage) == 0 {
		t.Fail()
	}
	out.Output()
}

func TestParseSenderIP(t *testing.T) {
	msg, err := msg.New("testdata/senderip.msg")
	if err != nil {
		t.Fail()
		return
	}

	out := msg.Format()
	if out.SenderAddress != "93.125.114.11" {
		t.Fail()
	}
	out.Output()
}

func TestParseMSG1(t *testing.T) {
	msg, err := msg.New("testdata/549970122456a12d8290cea3dd9c960f.msg")
	if err != nil {
		t.Fail()
		return
	}

	out := msg.Format()
	fmt.Println(out)
}

func TestParseEML1(t *testing.T) {
	msg, err := eml.New("testdata/6eabf11e48dbd66d451bdc03fc0d4913.eml")
	if err != nil {
		t.Fail()
		return
	}

	msg.Format().Output()
}
