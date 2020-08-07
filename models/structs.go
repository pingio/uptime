package models

import "net/http"

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
	Password string
	To       string
}
