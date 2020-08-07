package main

import (
	"encoding/json"
	"fmt"
	"github.com/pingio/uptime/models"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	// Our two structs that we use.
	var websites models.Websites
	var secret models.Secret

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
	responses := make(chan models.RequestResponse)

	numResponses := 0

	for _, website := range websites.Websites {
		for _, url := range website.Urls {
			numResponses++
			go models.Ping(url, responses)
		}
	}

	for i := 0; i < numResponses; i++ {
		response := <-responses
		if response.Response.StatusCode == 200 {
			log.Println(response.URL, "returned 200.")
		} else {
			models.Send(secret.Email, secret.Password, secret.To, response.URL, fmt.Sprintln(response.URL, "did not return a 200 OK, please investigate."))
		}
	}

}
