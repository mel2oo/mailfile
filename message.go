package mailfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
)

type Message struct {
	// 邮件头
	// Received段：路由信息，记录了邮件传递过程。
	Headers mail.Header `json:"-"`

	MessageID string `json:"message-id"`

	// 表示邮件建立的时间，既不是发送时间也不是接收时间，是邮件发送方创建邮件的时间。
	Date string `json:"date"`
	// 表示邮件的主题。
	Subject string `json:"subject"`

	// 发送者的ip地址
	SenderAddress string `json:"sender-address"`
	// 表示邮件的实际投递者（只能是一个），
	// 一般由收件方添加，邮件服务商在收到邮件后会将邮件会话里面的实际投递者与信头From字段标识的发件这进行比较，
	// 如不一致则在信头下方加入Sender字段标识邮件实际投递者，但这个字段也可由发件方决定的。
	Sender *mail.Address `json:"sender"`
	// 表示一个或多个邮件的作者，显示在正文的发件人。
	// 由发件方编辑，例如发垃圾的就会将此字段编辑成不存在的地址；发诈骗邮件的就会将此字段编辑成被冒充的邮件地址。
	From []*mail.Address `json:"from"`
	// 表示回复地址，由发件方编辑，希望收件人回复邮件时回复到指定的地址。
	// 一般情况下，如不额外添加Reply-to字段，收件人回复邮件时，将回复到原邮件From字段标识的地址。
	ReplyTo []*mail.Address `json:"reply-to"`

	// 表示邮件的接收地址。
	To []*mail.Address `json:"to"`
	// 表示抄送的邮件地址。
	Cc []*mail.Address `json:"cc"`
	// 表示密送的邮件地址。
	Bcc []*mail.Address `json:"bcc"`

	// 标识了邮件内容的格式
	ContentType string `json:"content-type"`

	// 邮件正文内容
	Body io.Reader `json:"-"`
	Html io.Reader `json:"-"`

	// 邮件正文中内嵌文件
	Embeddeds []Embedded `json:"embedded"`
	// 邮件附件
	Attachments []Attachment `json:"attachment"`
	// 邮件附件，子邮件类型
	SubMessage []*Message `json:"sub-message"`
}

type Attachment struct {
	Filename    string    `json:"filename"`
	ContentType string    `json:"content-type"`
	Data        io.Reader `json:"-"`
}

type Embedded struct {
	CID         string    `json:"cid"`
	ContentType string    `json:"content-type"`
	Data        io.Reader `json:"-"`
}

func GetSenderIP(headers mail.Header) (ip string, err error) {
	list, ok := headers["Received"]
	if !ok || len(list) == 0 {
		return ip, errors.New("received not found")
	}

	value := list[len(list)-1]
	left := strings.Index(value, "[")
	right := strings.Index(value, "]")
	if right-left < 7 {
		return ip, errors.New("address not found")
	}
	return value[left+1 : right], nil
}

func (m *Message) Output() {
	data := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(data)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(m); err != nil {
		fmt.Println("message output error")
	} else {
		fmt.Println(data.String())
	}
}
