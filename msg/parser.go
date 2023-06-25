package msg

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/mel2oo/mailfile"
	"github.com/richardlehane/mscfb"
)

func New(file string) (*Stream, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// MSCFB document reader
	doc, err := mscfb.New(f)
	if err != nil {
		return nil, err
	}

	// MSOX-MSG file stream reader
	return NewStream(doc)
}

func (s *Stream) Format() *mailfile.Message {
	msg := &mailfile.Message{}

	ParseProps(msg, s.UnpackData.props)
	ParseAttachment(msg, s.UnpackData.attachs)

	hdata, _ := ioutil.ReadAll(msg.Html)
	tdata, _ := ioutil.ReadAll(msg.Body)

	msg.Html = bytes.NewBuffer(hdata)
	msg.Body = bytes.NewBuffer(tdata)
	msg.Pwd = mailfile.ParsePasswd(hdata, tdata)
	return msg
}
