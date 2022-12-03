package msg

import (
	"os"

	"github.com/mel2oo/mailfile"
	"github.com/richardlehane/mscfb"
)

type MsOxMessage struct {
	stream  *Stream
	message *mailfile.Message
}

func New(file string) (*MsOxMessage, error) {
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
	stream, err := NewStream(doc)
	if err != nil {
		return nil, err
	}

	// MSOX-MSG stream data extract
	return &MsOxMessage{
		stream:  stream,
		message: Extract(stream),
	}, nil
}

func Extract(stream *Stream) *mailfile.Message {
	msg := &mailfile.Message{}

	ParseProps(msg, stream.props)
	ParseAttachment(msg, stream.attachs)

	return msg
}
