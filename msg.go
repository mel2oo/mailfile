package msgfile

import (
	"os"

	"github.com/richardlehane/mscfb"
)

type Message struct {
	props   MetaData
	subtag  []MetaData
	recips  []MetaData
	attachs []MetaData
}

func New(file string) (*Message, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	doc, err := mscfb.New(f)
	if err != nil {
		return nil, err
	}

	stream, err := readStream(doc)
	if err != nil {
		return nil, err
	}

	return stream.Extract(), nil
}
