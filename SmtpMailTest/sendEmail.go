package main

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
)

func main() {

	user := "test@rstc-service.com" //アカウント
	password := "UU&i19AQ"          //パスワード

	_from := mail.Address{"test@rstc-service.com", "test@rstc-service.com"}
	_to := mail.Address{"小林塁", "kobayashi.rui@gmail.com"}

	from := _from.String() //送信元のメールアドレス

	to := []string{
		_to.Address,
	}

	host := "smtp22.gmoserver.jp"
	addr := "smtp22.gmoserver.jp:587"

	msg := []byte("From: " + from + "\r\n" +
		"To: " + _to.String() + "\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Subject: Test mail\r\n\r\n" +
		"これはテストです！\r\n")

	auth := smtp.PlainAuth("", user, password, host)
	fmt.Println("ok auth")

	err := smtp.SendMail(addr, auth, from, to, msg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
