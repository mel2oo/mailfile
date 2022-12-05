package test

import (
	"fmt"
	"testing"

	"github.com/mel2oo/mailfile/eml"
	"github.com/mel2oo/mailfile/msg"
)

func TestParseMsg(t *testing.T) {
	msg, err := msg.New("testdata/sender_ip.msg")
	if err != nil {
		t.Fail()
		return
	}

	out := msg.Format()
	fmt.Println(out.Subject)
}

func TestParseEml(t *testing.T) {
	eml, err := eml.New("testdata/2.eml")
	if err != nil {
		t.Fail()
		return
	}

	out := eml.Format()
	fmt.Println(out.Subject)
}
