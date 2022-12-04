package mailfile

import (
	"io"
	"net/mail"
)

type Message struct {
	Headers mail.Header

	MessageID string   `json:"message-id"`
	Date      string   `json:"date"`
	Subject   string   `json:"subject"`
	Sender    string   `json:"sender"`
	From      []string `json:"from"`
	ReplyTo   []string `json:"reply-to"`
	To        []string `json:"to"`
	Cc        []string `json:"cc"`
	Bcc       []string `json:"bcc"`

	ContentType string
	Content     io.Reader

	Body        string
	Html        string
	Embeddeds   []Embedded
	Attachments []Attachment
	Child       []Message
}

type Attachment struct {
	Filename    string
	ContentType string
	Data        io.Reader
}

type Embedded struct {
	CID         string
	ContentType string
	Data        io.Reader
}

func (m *Message) Output() {
	// fmt.Println("MessageID ")
}
