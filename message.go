package mailfile

import (
	"io"
	"net/mail"
)

type Message struct {
	// 邮件头
	// Received段：路由信息，记录了邮件传递过程。
	Headers mail.Header

	MessageID string `json:"message-id"`

	// 表示邮件建立的时间，既不是发送时间也不是接收时间，是邮件发送方创建邮件的时间。
	Date string `json:"date"`
	// 表示邮件的主题。
	Subject string `json:"subject"`

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
	ContentType string

	// 邮件正文内容
	Body string
	Html string

	// 邮件正文中内嵌文件
	Embeddeds []Embedded
	// 邮件附件
	Attachments []Attachment
	// 邮件附件，子邮件类型
	SubMessage []Message
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
