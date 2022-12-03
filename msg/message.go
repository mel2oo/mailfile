package msg

import (
	"bytes"
	"net/mail"
	"strings"

	"github.com/mel2oo/mailfile"
)

func ParseProps(msg *mailfile.Message, m MetaData) {
	header, ok := m["TransportMessageHeaders"].(string)
	if ok {
		msg.Handers = Headers(header)
	}

	msg.MessageID, _ = m["InternetMessageId"].(string)

	msg.Date, ok = m["DeliverTime"].(string)
	if !ok {
		_, ok = msg.Handers["Date"]
		if ok {
			msg.Date = msg.Handers["Date"][0]
		}
	}

	msg.Subject, _ = m["Subject"].(string)

	_, ok = msg.Handers["Sender"]
	if ok {
		msg.Sender = msg.Handers["Sender"][0]
	}

	_, ok = msg.Handers["From"]
	if ok {
		msg.From = msg.Handers["From"]
	} else {
		from, ok := m["SenderRepresentingSmtpAddress"].(string)
		if ok {
			msg.From = []string{from}
		}
	}

	_, ok = msg.Handers["Reply-To"]
	if ok {
		msg.ReplyTo = msg.Handers["Reply-To"]
	} else {
		replyto, ok := m["ReplyRecipientNames"].(string)
		if ok {
			msg.ReplyTo = []string{replyto}
		}
	}

	_, ok = msg.Handers["TO"]
	if ok {
		msg.To = msg.Handers["TO"]
	} else {
		to, ok := m["DisplayTo"].(string)
		if !ok {
			to, _ = m["ReceivedRepresentingSmtpAddress"].(string)
		}
		msg.To = []string{to}
	}

	_, ok = msg.Handers["CC"]
	if ok {
		msg.Cc = msg.Handers["CC"]
	}

	_, ok = msg.Handers["BCC"]
	if ok {
		msg.Bcc = msg.Handers["BCC"]
	}

	msg.Body, _ = m["Body"].(string)

	html, ok := m["Html"].([]byte)
	if ok {
		msg.Html = string(html)
	}

	ctxtype, ok1 := msg.Handers["Content-Type"]
	ctxdata, ok2 := m["RtfCompressed"].([]uint8)
	if ok1 && ok2 {
		msg.ContentType = ctxtype[0]
		msg.Content = bytes.NewBuffer(ctxdata)
	}
}

func ParseRecipient(msg *mailfile.Message, datas []UnpackData) {
}

func ParseAttachment(msg *mailfile.Message, datas []UnpackData) {
	for _, data := range datas {

		filename, ok := data.props["AttachFilename"].(string)
		if ok {
			ctxtype, _ := data.props["AttachMimeTag"].(string)
			ctxdata, _ := data.props["AttachDataObject"].([]uint8)

			if len(ctxdata) == 0 && len(data.subtag) > 0 {
				ctxdata = data.subtag[0].props["RtfCompressed"].([]uint8)
			}

			msg.Attachments = append(msg.Attachments, mailfile.Attachment{
				Filename:    filename,
				ContentType: ctxtype,
				Data:        bytes.NewBuffer(ctxdata),
			})
			continue
		}

		cid, ok := data.props["AttachContentId"].(string)
		if ok {
			ctxtype, _ := data.props["AttachMimeTag"].(string)
			ctxdata, _ := data.props["AttachDataObject"].([]uint8)

			msg.Embeddeds = append(msg.Embeddeds, mailfile.Embedded{
				CID:         cid,
				ContentType: ctxtype,
				Data:        bytes.NewBuffer(ctxdata),
			})
		}
	}
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
