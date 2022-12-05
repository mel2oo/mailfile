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

func TestParseMSG(t *testing.T) {
	msg, err := msg.New("testdata/7d2e9038bb148560d795b383ddc7824b50f0916b4d2952262a1ba83a578e0453.msg")
	if err != nil {
		t.Fail()
		return
	}

	msg.Format().Output()
}

func TestParseEML(t *testing.T) {
	msg, err := eml.New("testdata/a854049f77696c2d7b4b5eee4707a9067d6fd94edd851023d3590829feccbd87.eml")
	if err != nil {
		t.Fail()
		return
	}

	msg.Format().Output()
}
