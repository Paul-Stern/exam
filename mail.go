// TODO: Export 'mail' to separate package
package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TODO: make local to mail package
var mailUser, mailHost, mailPort, mailPassword string

type EmailMessage struct {
	Date        time.Time
	From        string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

type Sender struct {
	auth smtp.Auth
}

func initMail() {
	mailUser = cfg.SMTP.User
	mailHost = cfg.SMTP.Host
	mailPort = cfg.SMTP.Port
	mailPassword = cfg.SMTP.Pswd

}

func NewSender() *Sender {
	auth := smtp.PlainAuth("", mailUser, mailPassword, mailHost)
	log.Printf("Auth: %+v", auth)
	return &Sender{auth}
}

func (s *Sender) Send(m *EmailMessage) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", mailHost, mailPort), s.auth, mailUser, m.To, m.ToBytes())
}

func NewMessage(s, b string) *EmailMessage {
	return &EmailMessage{
		Subject:     s,
		Date:        time.Now(),
		Body:        b,
		Attachments: make(map[string][]byte),
	}
}

func (m *EmailMessage) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

// TODO: Wrap long lines (MUST be no more than 998 and should be less 79)
// https://datatracker.ietf.org/doc/html/rfc2822#section-2.1.1
// TODO: Split into smaller methods
func (m *EmailMessage) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Date: %s\n", m.Date.Format(time.RFC822)))
	buf.WriteString(fmt.Sprintf("From: %s\n", m.From))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	buf.WriteString(
		fmt.Sprintf(
			"Subject: =?UTF-8?B?%s?=\n",
			base64.StdEncoding.EncodeToString([]byte(m.Subject)),
		),
	)
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n\r\n", http.DetectContentType([]byte(m.Body))))
	buf.WriteString(fmt.Sprintf("%s\r\n", m.Body))
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n\r\n", k))

			b64 := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b64, v)
			// Write attachment in chunks
			b := bytes.NewBuffer(b64)
			// Create chunk with gives size
			chunk := make([]byte, 77)
			for {
				// Read a chunk
				n, err := b.Read(chunk)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					log.Printf("Email attachment error: %s", err)
					break
				}
				// Write the chunk
				buf.Write(chunk[:n])
				buf.WriteString("\r\n")

			}
			buf.WriteString("\r\n")
		}
	}
	buf.WriteString(fmt.Sprintf("\r\n--%s--", boundary))

	return buf.Bytes()
}

func testEmail() (err error) {
	log.Print("Start testing mail services")
	initMail()
	c, err := smtp.Dial(fmt.Sprintf("%s:%s", cfg.SMTP.Host, cfg.SMTP.Port))
	if err != nil {
		log.Printf("Mail error: %s\nClient info: %+v", err, c)
		return
	}
	log.Printf("Successfully connected to smtp server: %+v", c)
	if cfg.SMTP.TestAddr == "" {
		log.Printf("No testing address. Test message won't be sent")
		log.Printf("Email service: success")
		return
	}
	m := NewMessage("Тест", "Это тестовое письмо")
	m.From = mailUser
	m.To = append(m.To, cfg.SMTP.TestAddr)
	// m.AttachFile("tmp/cert-2709233084.pdf")
	s := NewSender()
	err = s.Send(m)
	log.Print(err)
	return
}
