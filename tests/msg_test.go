package test

import (
	"fmt"
	"testing"

	"github.com/mel2oo/mailfile"
)

func TestParseMsg(t *testing.T) {
	msg, err := mailfile.New("testdata/complete.msg")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(msg)
}
