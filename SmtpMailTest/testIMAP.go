package main

import (
	"bytes"
	"log"
	"time"

	"github.com/emersion/go-imap/client"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("pop22.gmoserver.jp:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	user := "test@rstc-service.com" //アカウント
	password := "UU&i19AQ"          //パスワード
	if err := c.Login(user, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Write the message to a buffer
	var b bytes.Buffer
	b.WriteString("From: <hoge@test.com>\r\n")
	b.WriteString("To: <fuga@test.com>\r\n")
	b.WriteString("Subject: Hey there\r\n")
	b.WriteString("\r\n")
	b.WriteString("Hey <3")

	// Append it to INBOX, with two flags
	//flags := []string{imap.FlaggedFlag, "foobar"}
	flags := []string{}

	//write inbox
	//if err := c.Append("INBOX", flags, time.Now(), &b); err != nil {
	//	log.Fatal(err)
	//}

	if err := c.Append("Sent Messages", flags, time.Now(), &b); err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}
