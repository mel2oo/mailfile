package msg

import (
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
	msg.Pwd = mailfile.ParsePasswd(msg.Html, msg.Body)
	return msg
}
