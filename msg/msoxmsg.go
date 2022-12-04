package msg

import (
	"os"

	"github.com/mel2oo/mailfile"
	"github.com/richardlehane/mscfb"
)

type MsOxMessage struct {
	stream *Stream
	*mailfile.Message
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
		Message: Extract(stream.UnpackData),
	}, nil
}

func Extract(data UnpackData) *mailfile.Message {
	msg := &mailfile.Message{}

	ParseProps(msg, data.props)
	ParseAttachment(msg, data.attachs)

	return msg
}
