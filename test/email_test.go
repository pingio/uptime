package test

import (
	"encoding/json"
	"fmt"
	"github.com/pingio/uptime/models"
	"io/ioutil"
	"os"
	"testing"
)

func TestAbs(t *testing.T) {
	///////////////
	// Secret file.
	///////////////
	var secret models.Secret

	secretFile, err := os.Open("../secret.json")
	defer secretFile.Close()
	if err != nil {
		t.Error(err)
	}
	secretContent, err := ioutil.ReadAll(secretFile)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(secretContent, &secret)

	models.Send(secret.Email, secret.Password, secret.To, "test email", fmt.Sprintln("this is a test email."))
}
