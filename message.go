package mailfile

import (
	"io"
	"net/mail"
)

type Message struct {
	Handers mail.Header

	MessageID string
	Date      string
	Subject   string
	Sender    string
	From      []string
	ReplyTo   []string
	To        []string
	Cc        []string
	Bcc       []string

	ContentType string
	Content     io.Reader

	Body        string
	Html        string
	Embeddeds   []Embedded
	Attachments []Attachment
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
