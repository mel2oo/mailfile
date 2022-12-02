package msg

import (
	"os"

	"github.com/richardlehane/mscfb"
)

type MsOxMessage struct {
	*Message

	Attachments []string
	Recipients  []string
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

	return nil
}
