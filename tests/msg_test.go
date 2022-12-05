package test

import (
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
	if out.SenderAddress != "93.125.114.1" {
		t.Fail()
	}
	out.Output()
}

func TestParseEml(t *testing.T) {
	eml, err := eml.New("testdata/2.eml")
	if err != nil {
		t.Fail()
		return
	}

	eml.Format().Output()
}
