package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Websites struct {
	Websites []Website `json:"websites"`
}

type Website struct {
	Name string   `json:"name"`
	Urls []string `json:"urls"`
}

type RequestResponse struct {
	URL      string
	Response *http.Response
}

type Secret struct {
	Email    string
	To       string
	Password string
}

func main() {

	// Our two structs that we use.
	var websites Websites
	var secret Secret

	///////////////
	// Website file.
	///////////////
	websiteFile, err := os.Open("websites.json")
	defer websiteFile.Close()

	if err != nil {
		log.Fatalln(err)
	}

	// put the json content into the Websites struct.
	jsonContent, err := ioutil.ReadAll(websiteFile)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(jsonContent, &websites)

	///////////////
	// Secret file.
	///////////////
	secretFile, err := os.Open("secret.json")
	defer secretFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	secretContent, err := ioutil.ReadAll(secretFile)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(secretContent, &secret)

	///////////////////////////////////
	// Ping websites and send response.
	///////////////////////////////////
	responses := make(chan RequestResponse)

	numResponses := 0

	for _, website := range websites.Websites {
		for _, url := range website.Urls {
			numResponses++
			go ping(url, responses)
		}
	}

	for i := 0; i < numResponses; i++ {
		response := <-responses
		if response.Response.StatusCode == 200 {
			log.Println(response.URL, "returned 200.")
		} else {
			send(secret.Email, secret.Password, secret.To, response.URL, fmt.Sprintln(response.URL, "did not return a 200 OK, please investigate."))
		}
	}

}

func ping(website string, channel chan<- RequestResponse) {
	response, err := http.Get(website)

	if err != nil {
		log.Println(website, "could not get connection.")
	}
	reqresp := RequestResponse{website, response}
	channel <- reqresp
}

func send(email, password, sendTo, url, body string) {
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
