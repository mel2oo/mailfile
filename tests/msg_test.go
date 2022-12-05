package test

import (
	"fmt"
	"io"
	"testing"

	"github.com/mel2oo/mailfile/eml"
	"github.com/mel2oo/mailfile/msg"
	"github.com/stretchr/testify/assert"
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
	msg, err := eml.New("testdata/476ae97d5536c2712f455f633c0c1ff7.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, len(res.Embeddeds), 5)
}

func TestParseEML2(t *testing.T) {
	msg, err := eml.New("testdata/6eabf11e48dbd66d451bdc03fc0d4913.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()

	for _, m := range res.SubMessage {
		assert.Equal(t, m.Subject, "Payment Receipt")

		for _, a := range m.Attachments {
			assert.Equal(t, a.Filename, "parcel2go.com.html")

			data, _ := io.ReadAll(a.Data)
			assert.Greater(t, len(data), 0)
		}
	}
}

func TestParseEML3(t *testing.T) {
	msg, err := eml.New("testdata/927e94db23827c4247e112f04ff80769.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()

	for _, f := range res.From {
		assert.Equal(t, f.Name, "RAFAEL ANTONIO CASAS TAPIAS")
		assert.Equal(t, f.Address, "rafael.casasta@unaula.edu.co")
	}

	for _, a := range res.Attachments {
		assert.Equal(t, a.Filename, "8513607623965074005428406083054411951091040465890343326663102629660275767.tgz")
	}
}

func TestParseEML4(t *testing.T) {
	msg, err := eml.New("testdata/db84a1ca6bd634d671e39908bc3f3e0e.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, res.Attachments[0].Filename, "KYC2633.html")
	assert.Equal(t, res.SenderAddress, "71.168.222.19")
}
