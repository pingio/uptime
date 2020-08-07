package models

import (
	"log"
	"net/http"
	"net/smtp"
)

func Ping(website string, channel chan<- RequestResponse) {
	response, err := http.Get(website)

	if err != nil {
		log.Println(website, "could not get connection.")
	}
	reqresp := RequestResponse{website, response}
	channel <- reqresp
}

func Send(email, password, sendTo, url, body string) {
	from := email
	pass := password
	to := sendTo

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + url + " error\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
