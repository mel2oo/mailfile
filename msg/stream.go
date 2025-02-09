package msg

import (
	"io"
	"strings"

	"github.com/richardlehane/mscfb"
)

type Stream struct {
	// unpack
	UnpackData
	// origin
	origin msoxstream
}

func NewStream(doc *mscfb.Reader) (*Stream, error) {
	stream := &Stream{
		origin: msoxstream{
			props:   make([]*mscfb.File, 0),
			subtag:  make(map[string]*msoxstream),
			recips:  make(map[string]*msoxstream),
			attachs: make(map[string]*msoxstream),
		},
	}

	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if strings.Contains(entry.Name, "__substg1.0_") {
			stream.origin.setEntry(entry.Path, entry)
		}
	}

	stream.UnpackData = stream.origin.extract()

	return stream, nil
}

type msoxstream struct {
	props   []*mscfb.File
	subtag  map[string]*msoxstream
	recips  map[string]*msoxstream
	attachs map[string]*msoxstream
}

func (s *msoxstream) setEntry(keys []string, entry *mscfb.File) {
	if len(keys) == 0 {
		s.props = append(s.props, entry)
	} else {

		if strings.Contains(keys[0], "__substg1.0_") {
			if s.subtag[keys[0]] == nil {
				s.subtag[keys[0]] = &msoxstream{
					props:   make([]*mscfb.File, 0),
					subtag:  make(map[string]*msoxstream),
					recips:  make(map[string]*msoxstream),
					attachs: make(map[string]*msoxstream),
				}
			}
			s.subtag[keys[0]].setEntry(keys[1:], entry)
		}

		if strings.Contains(keys[0], "__attach_") {
			if s.attachs[keys[0]] == nil {
				s.attachs[keys[0]] = &msoxstream{
					props:   make([]*mscfb.File, 0),
					subtag:  make(map[string]*msoxstream),
					recips:  make(map[string]*msoxstream),
					attachs: make(map[string]*msoxstream),
				}
			}
			s.attachs[keys[0]].setEntry(keys[1:], entry)
		}

		if strings.Contains(keys[0], "__recip_") {
			if s.recips[keys[0]] == nil {
				s.recips[keys[0]] = &msoxstream{
					props:   make([]*mscfb.File, 0),
					subtag:  make(map[string]*msoxstream),
					recips:  make(map[string]*msoxstream),
					attachs: make(map[string]*msoxstream),
				}
			}
			s.recips[keys[0]].setEntry(keys[1:], entry)
		}
	}
}

type MetaData map[string]interface{}

type UnpackData struct {
	props   MetaData
	subtag  []UnpackData
	recips  []UnpackData
	attachs []UnpackData
}

func (m *msoxstream) extract() UnpackData {
	up := UnpackData{
		props:   m.unpack(),
		subtag:  make([]UnpackData, 0),
		recips:  make([]UnpackData, 0),
		attachs: make([]UnpackData, 0),
	}

	for _, maps := range m.subtag {
		up.subtag = append(up.subtag, maps.extract())
	}

	for _, maps := range m.recips {
		up.recips = append(up.recips, maps.extract())
	}

	for _, maps := range m.attachs {
		up.attachs = append(up.attachs, maps.extract())
	}

	return up
}

func (m *msoxstream) unpack() MetaData {
	var (
		metadata              = make(MetaData)
		directory_name_filter = "__substg1.0_"
	)

	for _, entry := range m.props {
		if entry == nil {
			continue
		}

		if !strings.Contains(entry.Name, directory_name_filter) {
			continue
		}

		property_name, property_type := m.PropsNameType(entry)
		if len(property_name) == 0 {
			continue
		}

		data, err := io.ReadAll(entry)
		if err != nil {
			continue
		}

		if property_name == "AttachDataObject" {
			metadata[property_name] = data
		} else {
			metadata[property_name] = GetDataValue(property_type, data)
		}
	}

	return metadata
}

func (m *msoxstream) PropsNameType(entry *mscfb.File) (property_name, property_type string) {
	if strings.Contains(entry.Name, "__substg1.0_") {
		namid := "0x" + strings.ReplaceAll(entry.Name, "__substg1.0_", "")[0:4]
		property_type = "0x" + strings.ReplaceAll(entry.Name, "__substg1.0_", "")[4:8]
		props := PROPS_ID_MAP[namid]
		if property_type != "0x0000" {
			return props["name"], property_type
		}
		return props["name"], props["data_type"]
	}

	return
}
