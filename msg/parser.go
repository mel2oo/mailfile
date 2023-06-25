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

	var hdata, tdata []byte

	if msg.Html != nil {
		hdata, _ = ioutil.ReadAll(msg.Html)
	}
	if msg.Body != nil {
		tdata, _ = ioutil.ReadAll(msg.Body)
	}

	if len(hdata) > 0 {
		msg.Html = bytes.NewBuffer(hdata)
	}
	if len(tdata) > 0 {
		msg.Body = bytes.NewBuffer(tdata)
	}

	msg.Pwd = mailfile.ParsePasswd(hdata, tdata)
	return msg
}
