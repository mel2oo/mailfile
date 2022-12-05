package eml

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"

	"github.com/mel2oo/mailfile"
)

func New(file string) (*Message, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return ParseMessage(fi)
}

// ParseMessage parses and returns a Message from an io.Reader
// containing the raw text of an email message.
// (If the raw email is a string or []byte, use strings.NewReader()
// or bytes.NewReader() to create a reader.)
// Any "quoted-printable" or "base64" encoded bodies will be decoded.
func ParseMessage(r io.Reader) (*Message, error) {
	msg, err := mail.ReadMessage(&leftTrimReader{r: bufioReader(r)})
	if err != nil {
		return nil, err
	}
	// decode any Q-encoded values
	for _, values := range msg.Header {
		for idx, val := range values {
			values[idx] = decodeRFC2047(val)
		}
	}
	return parseMessageWithHeader(Header(msg.Header), msg.Body)
}

// parseMessageWithHeader parses and returns a Message from an already filled
// Header, and an io.Reader containing the raw text of the body/payload.
// (If the raw body is a string or []byte, use strings.NewReader()
// or bytes.NewReader() to create a reader.)
// Any "quoted-printable" or "base64" encoded bodies will be decoded.
func parseMessageWithHeader(headers Header, bodyReader io.Reader) (*Message, error) {

	bufferedReader := contentReader(headers, bodyReader)

	var err error
	var mediaType string
	var mediaTypeParams map[string]string
	var preamble []byte
	var epilogue []byte
	var body []byte
	var parts []*Message
	var subMessage *Message

	if contentType := headers.Get("Content-Type"); len(contentType) > 0 {
		mediaType, mediaTypeParams, err = mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}
	} // Lack of contentType is not a problem

	// Can only have one of the following: Parts, SubMessage, or Body
	if strings.HasPrefix(mediaType, "multipart") {
		boundary := mediaTypeParams["boundary"]
		preamble, err = readPreamble(bufferedReader, boundary)
		if err == nil {
			parts, err = readParts(bufferedReader, boundary)
			if err == nil {
				epilogue, err = readEpilogue(bufferedReader)
			}
		}

	} else if strings.HasPrefix(mediaType, "message") {
		subMessage, err = ParseMessage(bufferedReader)

	} else {
		body, err = ioutil.ReadAll(bufferedReader)
	}
	if err != nil {
		return nil, err
	}

	return &Message{
		Header:     headers,
		Preamble:   preamble,
		Epilogue:   epilogue,
		Body:       body,
		SubMessage: subMessage,
		Parts:      parts,
	}, nil
}

// readParts parses out the parts of a multipart body, including the preamble and epilogue.
func readParts(bodyReader io.Reader, boundary string) ([]*Message, error) {

	parts := make([]*Message, 0, 1)
	multipartReader := multipart.NewReader(bodyReader, boundary)

	for part, partErr := multipartReader.NextPart(); partErr != io.EOF; part, partErr = multipartReader.NextPart() {
		if partErr != nil && partErr != io.EOF {
			return []*Message{}, partErr
		}
		newEmailPart, msgErr := parseMessageWithHeader(Header(part.Header), part)
		part.Close()
		if msgErr != nil {
			return []*Message{}, msgErr
		}
		parts = append(parts, newEmailPart)
	}
	return parts, nil
}

// readEpilogue ...
func readEpilogue(r io.Reader) ([]byte, error) {
	epilogue, err := ioutil.ReadAll(r)
	for len(epilogue) > 0 && isASCIISpace(epilogue[len(epilogue)-1]) {
		epilogue = epilogue[:len(epilogue)-1]
	}
	if len(epilogue) > 0 {
		return epilogue, err
	}
	return nil, err
}

// readPreamble ...
func readPreamble(r *bufio.Reader, boundary string) ([]byte, error) {
	preamble, err := ioutil.ReadAll(&preambleReader{r: r, boundary: []byte("--" + boundary)})
	if len(preamble) > 0 {
		return preamble, err
	}
	return nil, err
}

// preambleReader ...
type preambleReader struct {
	r        *bufio.Reader
	boundary []byte
}

// Read ...
func (r *preambleReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Peek and read up to the --boundary, then EOF
	peek, err := r.r.Peek(len(p))
	if err != nil && err != io.EOF {
		return 0, err
	}

	idx := bytes.Index(peek, r.boundary)

	if idx < 0 {
		// Couldn't find the boundary, so read all the bytes we can,
		// but leave room for a new-line + boundary that got cut in half by the buffer,
		// that way it can be matched against on the next read
		return r.r.Read(p[:max(1, len(peek)-(len(r.boundary)+2))])
	}

	// Account for possible new-line / whitespace at start of the boundary, which shouldn't be removed
	for idx > 0 && isASCIISpace(peek[idx-1]) {
		idx--
	}

	if idx == 0 {
		// The boundary (or new-line + boundary) is at the start of the reader, so there is no preamble
		return 0, io.EOF
	}

	n, err := r.r.Read(p[:idx])
	if err != nil && err != io.EOF {
		return n, err
	}
	return n, io.EOF
}

// contentReader ...
func contentReader(headers Header, bodyReader io.Reader) *bufio.Reader {
	if headers.Get("Content-Transfer-Encoding") == "quoted-printable" {
		headers.Del("Content-Transfer-Encoding")
		return bufioReader(quotedprintable.NewReader(bodyReader))
	}
	if headers.Get("Content-Transfer-Encoding") == "base64" {
		headers.Del("Content-Transfer-Encoding")
		return bufioReader(base64.NewDecoder(base64.StdEncoding, bodyReader))
	}
	return bufioReader(bodyReader)
}

// decodeRFC2047 ...
func decodeRFC2047(s string) string {
	// GO 1.5 does not decode headers, but this may change in future releases...
	decoded, err := (&mime.WordDecoder{}).DecodeHeader(s)
	if err != nil || len(decoded) == 0 {
		return s
	}
	return decoded
}

func (m *Message) Format() *mailfile.Message {
	var msg mailfile.Message
	msg.Headers = mail.Header(m.Header)
	msg.MessageID = m.Header.Get("Message-Id")
	msg.Date = m.Header.Get("Date")
	msg.Subject = m.Header.Subject()
	msg.ContentType = m.Header.Get("Content-Type")

	msg.SenderAddress, _ = mailfile.GetSenderIP(msg.Headers)
	msg.Sender, _ = mail.ParseAddress(m.Header.Get("Sender"))
	msg.From, _ = mail.ParseAddressList(m.Header.From())
	msg.ReplyTo, _ = mail.ParseAddressList(m.Header.Get("Reply-To"))

	for _, tstr := range m.Header.To() {
		to, err := mail.ParseAddress(tstr)
		if err == nil {
			msg.To = append(msg.To, to)
		}
	}

	for _, cstr := range m.Header.Cc() {
		cc, err := mail.ParseAddress(cstr)
		if err == nil {
			msg.Cc = append(msg.Cc, cc)
		}
	}

	for _, cstr := range m.Header.Bcc() {
		bcc, err := mail.ParseAddress(cstr)
		if err == nil {
			msg.Bcc = append(msg.Bcc, bcc)
		}
	}

	if m.HasBody() {
		msg.Body = bytes.NewBuffer(m.Body)
	}

	for index, part := range m.Parts {
		str, desc, err := part.Header.ContentDisposition()
		if err != nil {
			for _, partdata := range part.Parts {
				mime, _, err := partdata.Header.ContentType()
				if err != nil {
					continue
				}

				switch mime {
				case "text/plain":
					msg.Body = bytes.NewBuffer(partdata.Body)
				case "text/html":
					msg.Html = bytes.NewBuffer(partdata.Body)
				}
			}
			continue
		}

		if str == "attachment" {
			name, ok := desc["filename"]
			if !ok {
				name = fmt.Sprintf("%s_%d", str, index)
			}

			msg.Attachments = append(msg.Attachments, mailfile.Attachment{
				Filename:    name,
				ContentType: part.Header.Get("Content-Type"),
				Data:        bytes.NewBuffer(part.Body),
			})
		} else {
			cid, ok := desc["CID"]
			if !ok {
				cid = fmt.Sprintf("%s_%d", str, index)
			}

			msg.Embeddeds = append(msg.Embeddeds, mailfile.Embedded{
				CID:         cid,
				ContentType: part.Header.Get("Content-Type"),
				Data:        bytes.NewBuffer(part.Body),
			})
		}
	}

	return &msg
}
