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

	_, ok = msg.Headers["From"]
	if ok {
		for _, str1 := range msg.Headers["From"] {
			addrs, err := mail.ParseAddressList(str1)
			if err == nil {
				msg.From = append(msg.From, addrs...)
			}
		}
	} else {
		fromlist, ok := m["SenderRepresentingSmtpAddress"].(string)
		if ok {
			msg.From, _ = mail.ParseAddressList(fromlist)
		}
	}

	_, ok = msg.Headers["Sender"]
	if ok {
		msg.Sender, _ = mail.ParseAddress(msg.Headers["Sender"][0])
	} else {
		if len(msg.From) > 0 {
			msg.Sender = msg.From[0]
		}
	}

	_, ok = msg.Headers["Reply-To"]
	if ok {
		for _, str1 := range msg.Headers["Reply-To"] {
			addrs, err := mail.ParseAddressList(str1)
			if err == nil {
				msg.ReplyTo = append(msg.ReplyTo, addrs...)
			}
		}
	} else {
		replytolist, ok := m["ReplyRecipientNames"].(string)
		if ok {
			msg.ReplyTo, _ = mail.ParseAddressList(replytolist)
		}
	}

	to1, ok1 := msg.Headers["To"]
	to2, ok2 := msg.Headers["DisplayTo"]
	if ok1 || ok2 {
		for _, str1 := range append(to1, to2...) {
			addrs, err := mail.ParseAddressList(str1)
			if err == nil {
				msg.To = append(msg.To, addrs...)
			}
		}
	} else {
		tolist, ok := m["ReceivedRepresentingSmtpAddress"].(string)
		if ok {
			msg.To, _ = mail.ParseAddressList(tolist)
		}
	}

	_, ok = msg.Headers["CC"]
	if ok {
		for _, str1 := range msg.Headers["CC"] {
			addrs, err := mail.ParseAddressList(str1)
			if err == nil {
				msg.Cc = append(msg.Cc, addrs...)
			}
		}
	}

	_, ok = msg.Headers["BCC"]
	if ok {
		for _, str1 := range msg.Headers["BCC"] {
			addrs, err := mail.ParseAddressList(str1)
			if err == nil {
				msg.Bcc = append(msg.Bcc, addrs...)
			}
		}
	}

	msg.Body, _ = m["Body"].(string)

	html, ok := m["Html"].([]byte)
	if ok {
		msg.Html = string(html)
	}

	ctxtype, ok1 := msg.Headers["Content-Type"]
	if ok1 {
		msg.ContentType = ctxtype[0]
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
					msg.SubMessage = append(msg.SubMessage, msgfile)
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
