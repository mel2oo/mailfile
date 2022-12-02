package msg

import (
	"net/mail"
	"strings"
)

type Message struct {
	Handers      mail.Header
	MessageID    string
	Subject      string
	CreatedDate  string
	ReceivedDate string
	Date         string
	From         string
	To           string
	CC           string
	BCC          string
	ReplyTo      string
	Body         string
	Html         []byte
}

func ParseProps(m MetaData) Message {
	var msg Message

	header, ok := m["TransportMessageHeaders"].(string)
	if ok {
		msg.Handers = Headers(header)
	}

	msg.MessageID, _ = m["InternetMessageId"].(string)
	msg.Subject, _ = m["Subject"].(string)
	msg.CreatedDate, _ = m["CreationTime"].(string)
	msg.ReceivedDate, _ = m["ReceiptTime"].(string)
	msg.CC, _ = m["DisplayCc"].(string)
	msg.BCC, _ = m["DisplayBcc"].(string)
	msg.Body, _ = m["Body"].(string)
	msg.Html, _ = m["Html"].([]byte)

	msg.Date, ok = m["DeliverTime"].(string)
	if !ok {
		msg.Date, _ = m["Date"].(string)
	}

	msg.From, ok = m["From"].(string)
	if !ok {
		msg.From, _ = m["SenderRepresentingSmtpAddress"].(string)
	}

	msg.To, ok = m["TO"].(string)
	if !ok {
		msg.To, ok = m["DisplayTo"].(string)
		if !ok {
			msg.To, _ = m["ReceivedRepresentingSmtpAddress"].(string)
		}
	}

	msg.ReplyTo, ok = m["Reply-To"].(string)
	if !ok {
		msg.ReplyTo, _ = m["ReplyRecipientNames"].(string)
	}

	return msg
}

type Recipient struct {
}

func ParseRecipient(m []MetaData) []Recipient {
	return nil
}

type Attachment struct {
}

func ParseAttachment(m []MetaData) []Attachment {
	return nil
}

func Headers(hstr string) mail.Header {
	var (
		headers = make(mail.Header)
		key     string
		val     string
	)

	list := strings.Split(hstr, "\r\n")
	for _, s := range list {
		if strings.Contains(s, ": ") {
			index := strings.Index(s, ": ")

			if len(key) > 0 {
				if _, ok := headers[key]; !ok {
					headers[key] = make([]string, 0)
				}
				headers[key] = append(headers[key], val)
			}

			key = s[:index]
			val = s[index+2:]
		} else {
			val += s
		}
	}

	return headers
}
