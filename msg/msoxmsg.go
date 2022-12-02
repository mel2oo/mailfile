package msg

import (
	"os"

	"github.com/richardlehane/mscfb"
)

type MsOxMessage struct {
	Message

	Recipients  []Recipient
	Attachments []Attachment
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
	return Extract(stream), nil
}

func Extract(stream *Stream) *MsOxMessage {
	return &MsOxMessage{
		Message:     ParseProps(stream.props),
		Recipients:  ParseRecipient(stream.recips),
		Attachments: ParseAttachment(stream.attachs),
	}
}
