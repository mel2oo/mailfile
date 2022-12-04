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
		msg.Headers = Headers(header)
	}

	msg.MessageID, _ = m["InternetMessageId"].(string)

	msg.Date, ok = m["DeliverTime"].(string)
	if !ok {
		_, ok = msg.Headers["Date"]
		if ok {
			msg.Date = msg.Headers["Date"][0]
		}
	}

	msg.Subject, _ = m["Subject"].(string)

	_, ok = msg.Headers["Sender"]
	if ok {
		msg.Sender = msg.Headers["Sender"][0]
	}

	_, ok = msg.Headers["From"]
	if ok {
		msg.From = msg.Headers["From"]
	} else {
		from, ok := m["SenderRepresentingSmtpAddress"].(string)
		if ok {
			msg.From = []string{from}
		}
	}

	_, ok = msg.Headers["Reply-To"]
	if ok {
		msg.ReplyTo = msg.Headers["Reply-To"]
	} else {
		replyto, ok := m["ReplyRecipientNames"].(string)
		if ok {
			msg.ReplyTo = []string{replyto}
		}
	}

	_, ok = msg.Headers["TO"]
	if ok {
		msg.To = msg.Headers["TO"]
	} else {
		to, ok := m["DisplayTo"].(string)
		if !ok {
			to, _ = m["ReceivedRepresentingSmtpAddress"].(string)
		}
		msg.To = []string{to}
	}

	_, ok = msg.Headers["CC"]
	if ok {
		msg.Cc = msg.Headers["CC"]
	}

	_, ok = msg.Headers["BCC"]
	if ok {
		msg.Bcc = msg.Headers["BCC"]
	}

	msg.Body, _ = m["Body"].(string)

	html, ok := m["Html"].([]byte)
	if ok {
		msg.Html = string(html)
	}

	ctxtype, ok1 := msg.Headers["Content-Type"]
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
			if len(ctxdata) > 0 {
				msg.Attachments = append(msg.Attachments, mailfile.Attachment{
					Filename:    filename,
					ContentType: ctxtype,
					Data:        bytes.NewBuffer(ctxdata),
				})
			}

			if len(data.subtag) > 0 {
				for _, subdata := range data.subtag {
					var msgfile mailfile.Message
					ParseProps(&msgfile, subdata.props)
					ParseAttachment(&msgfile, subdata.attachs)
					msg.Child = append(msg.Child, msgfile)
				}
			}

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
