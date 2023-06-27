package test

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/mel2oo/mailfile/eml"
	"github.com/mel2oo/mailfile/msg"
	"github.com/stretchr/testify/assert"
)

func TestParseMSG1(t *testing.T) {
	msg, err := msg.New("testdata/1b098cd4bc21836a74d12ec519bd0e8c.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), "From: Frances Evensen")
	assert.Equal(t, res.SubMessage[0].Subject, "Frances Evensen shared \"SecureSave RFP\" with you.")
	assert.Equal(t, len(res.SubMessage[0].Embeddeds), 4)
}

func TestParseMSG2(t *testing.T) {
	msg, err := msg.New("testdata/0bb5983192375432403c74cf2d68ee67.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	html, _ := io.ReadAll(res.Html)
	assert.Contains(t, string(html), "<html><head>\r\n<meta http-equiv=")
	assert.Equal(t, len(res.Embeddeds), 2)
}

func TestParseMSG3(t *testing.T) {
	msg, err := msg.New("testdata/5499732e4b2d8f6da3f053e086ee479f.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, res.SenderAddress, "127.0.0.1")
	assert.Equal(t, res.Attachments[0].Filename, "▶🔘─────.htm")
}

func TestParseMSG4(t *testing.T) {
	msg, err := msg.New("testdata/7378473901a31ba720324e40d7fb1b3a.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, res.Attachments[0].Filename, "NIT SUSPENDIDO DETALLES DIAN.pdf")
}

func TestParseMSG5(t *testing.T) {
	msg, err := msg.New("testdata/549970122456a12d8290cea3dd9c960f.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, res.SubMessage[0].Subject, "message sent from (405)-3633914")
	assert.Equal(t, res.SubMessage[0].Attachments[0].Filename, "SKM59469.htm")
}

func TestParseMSG6(t *testing.T) {
	msg, err := msg.New("testdata/b9bd32895692dc99fa046b0655ef170f.msg")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, len(res.Attachments), 3)
	assert.Equal(t, len(res.SubMessage), 6)
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

func TestParseEML5(t *testing.T) {
	msg, err := eml.New("testdata/f21fda978107a91b5b6e3f3b50f533f9.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	body, _ := io.ReadAll(res.Body)
	html, _ := io.ReadAll(res.Html)
	assert.Contains(t, string(body), "You recently received a")
	assert.Contains(t, string(html), "<html xmlns:o=\"urn:schemas-")
}

func TestParseEML6(t *testing.T) {
	msg, err := eml.New("testdata/f93a4468e029ca3c81c800c80028d9b2.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	assert.Equal(t, len(res.SubMessage[0].Headers), 3)
	assert.Equal(t, res.SubMessage[1].Subject, "Fwd: Contract & Deposite//Revised Order")
}

func TestParseEML7(t *testing.T) {
	msg, err := eml.New("testdata/字符编码测试.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	t.Log(res.From[0].Name)
}

func TestDecode(t *testing.T) {
	a := "1cXB+g=="
	sDec, err := base64.StdEncoding.DecodeString(a)
	fmt.Println(string(sDec), err)

	t.Log(eml.IsGBK(sDec))

	d, err := eml.GbkToUtf8(sDec)
	if err != nil {
		t.Log(err)
	}

	t.Log(string(d))
}

func TestDecodeAllPasswd(t *testing.T) {
	filepath.Walk("passwd", func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		msg, err := eml.New(path)
		if err != nil {
			t.Fail()
			return nil
		}

		res := msg.Format()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("file: %s\n", path)
		fmt.Printf("body: %s\n", string(body))
		fmt.Printf("password: %v\n", res.Pwd)
		fmt.Printf("------------------------------------------------\n")
		return nil
	})
}
func TestDecodeOnePasswd(t *testing.T) {
	msg, err := eml.New("passwd/2看看密码是多少.eml")
	if err != nil {
		t.Fail()
		return
	}

	res := msg.Format()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	t.Log(res.From[0].Name)

}
