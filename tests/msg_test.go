package test

import (
	"fmt"
	"testing"

	"github.com/mel2oo/mailfile/msg"
)

func TestParseMsg(t *testing.T) {
	msg, err := msg.New("testdata/complete.msg")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(msg)
}
