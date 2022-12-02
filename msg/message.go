package msg

type Message struct {
	MessageID string
}

func ParseProps(m *MetaData) *Message {
	return nil
}

type Attachment struct {
}

func ParseAttach(m *MetaData) *Attachment {
	return nil
}

type Recipient struct {
}

func ParseRecipient(m *MetaData) *Recipient {
	return nil
}
