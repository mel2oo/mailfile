package msgfile

import (
	"strings"

	"github.com/richardlehane/mscfb"
)

type Stream struct {
	props   []*mscfb.File
	subtag  map[string]*Stream
	recips  map[string]*Stream
	attachs map[string]*Stream
}

func readStream(doc *mscfb.Reader) (*Stream, error) {
	s := newStream()

	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if strings.Contains(entry.Name, "__substg1.0_") {
			s.setEntry(entry.Path, entry)
		}
	}

	return s, nil
}

func newStream() *Stream {
	return &Stream{
		props:   make([]*mscfb.File, 0),
		subtag:  make(map[string]*Stream),
		recips:  make(map[string]*Stream),
		attachs: make(map[string]*Stream),
	}
}

func (s *Stream) setEntry(keys []string, entry *mscfb.File) {
	if len(keys) == 0 {
		s.props = append(s.props, entry)
	} else {

		if strings.Contains(keys[0], "__substg1.0_") {
			if s.subtag[keys[0]] == nil {
				s.subtag[keys[0]] = newStream()
			}
			s.subtag[keys[0]].setEntry(keys[1:], entry)
		}

		if strings.Contains(keys[0], "__attach_") {
			if s.attachs[keys[0]] == nil {
				s.attachs[keys[0]] = newStream()
			}
			s.attachs[keys[0]].setEntry(keys[1:], entry)
		}

		if strings.Contains(keys[0], "__recip_") {
			if s.recips[keys[0]] == nil {
				s.recips[keys[0]] = newStream()
			}
			s.recips[keys[0]].setEntry(keys[1:], entry)
		}
	}
}
