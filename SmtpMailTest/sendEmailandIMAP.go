package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/emersion/go-imap/client"
)

type Email struct {
	User     string
	Password string
	From     string
	Host     string
	SmtpAddr string
	ImapAddr string
	Auth     smtp.Auth
	Client   *client.Client
}

func (m *Email) writeString(b *bytes.Buffer, s string) *bytes.Buffer {
	_, err := b.WriteString(s)
	if err != nil {
		fmt.Print(err.Error())
	}
	return b
}

// サブジェクトを MIME エンコードする
func (m *Email) encodeSubject(subject string) string {
	// UTF8 文字列を指定文字数で分割する
	b := bytes.NewBuffer([]byte(""))
	strs := []string{}
	length := 13
	for k, c := range strings.Split(subject, "") {
		b.WriteString(c)
		if k%length == length-1 {
			strs = append(strs, b.String())
			b.Reset()
		}
	}
	if b.Len() > 0 {
		strs = append(strs, b.String())
	}
	// MIME エンコードする
	b2 := bytes.NewBuffer([]byte(""))
	b2.WriteString("Subject:")
	for _, line := range strs {
		b2.WriteString(" =?utf-8?B?")
		b2.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		b2.WriteString("?=\r\n")
	}
	return b2.String()
}

// 本文を 76 バイト毎に CRLF を挿入して返す
func (m *Email) encodeBody(body string) string {
	b := bytes.NewBufferString(body)
	s := base64.StdEncoding.EncodeToString(b.Bytes())
	b2 := bytes.NewBuffer([]byte(""))
	for k, c := range strings.Split(s, "") {
		b2.WriteString(c)
		if k%76 == 75 {
			b2.WriteString("\r\n")
		}
	}
	return b2.String()

}

func (e *Email) EmailAuth() error {
	e.Auth = smtp.PlainAuth("", e.User, e.Password, e.Host)
	return nil
}

func (e *Email) SetImap() error {
	c, err := client.DialTLS(e.ImapAddr, nil)
	if err != nil {
		return err
	}
	e.Client = c
	return nil
}

func (e *Email) Setup(user string, password string, fromName string, fromMail string, host string, smtpaddr string, imapaddr string) error {
	e.User = user
	e.Password = password
	_form := mail.Address{fromName, fromMail}
	e.From = _form.String()
	e.Host = host
	e.SmtpAddr = smtpaddr
	e.ImapAddr = imapaddr
	err := e.EmailAuth()
	if err != nil {
		return err
	}

	err = e.SetImap()
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) Send(to string, to_mail string, subject string, message string) error {
	//var msg []byte
	_to := mail.Address{to, to_mail}
	msg := bytes.NewBuffer([]byte(""))
	msg = e.writeString(msg, "From: "+e.From+"\r\n")
	msg = e.writeString(msg, "To: "+_to.String()+"\r\n")
	msg = e.writeString(msg, "MIME-Version: 1.0\r\n")
	msg = e.writeString(msg, "Content-Type: text/plain; charset=utf-8\r\n")
	msg = e.writeString(msg, "Content-Transfer-Encoding: base64\r\n")
	msg = e.writeString(msg, "Subject: "+subject+"\r\n")
	msg = e.writeString(msg, "\r\n")
	msg = e.writeString(msg, e.encodeBody(message))

	err := smtp.SendMail(e.SmtpAddr, e.Auth, e.From, []string{_to.String()}, msg.Bytes())
	if err != nil {
		return err
	}
	err = e.SetSendMail(msg)

	if err != nil {
		return err
	}
	return nil
}

// 送信済みに適用
func (e *Email) SetSendMail(message *bytes.Buffer) error {

	defer e.Client.Logout()

	if err := e.Client.Login(e.User, e.Password); err != nil {
		return err
	}
	log.Println("Logged in")

	flags := []string{}

	if err := e.Client.Append("Sent", flags, time.Now(), message); err != nil {
		return err
	}

	return nil
}

func main() {

	var emailController Email

	err := emailController.Setup(
		"test@rstc-service.com",
		"UU&i19AQ",
		"RSTC Test",
		"test@rstc-service.com",
		"smtp22.gmoserver.jp",
		"smtp22.gmoserver.jp:587",
		"pop22.gmoserver.jp:993")

	if err != nil {
		log.Fatal(err)
		return
	}

	err = emailController.Send("", "kobayashi.rui@gmail.com", "これはテスト送信です", "Hello This is Test. \n テスト送信です.")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Email sent successfully")
}
