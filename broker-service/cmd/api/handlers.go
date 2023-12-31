package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse {
		Error: false,
		Message: "Break the Broker",
	}

	app.writeJson(w, http.StatusOK, payload)
	
	fmt.Printf("Payload call")
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var rp RequestPayload

	err := app.readJson(w, r, &rp)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch rp.Action {
	case "auth":
		app.authenticate(w, rp.Auth)
	case "log":
		app.logItem(w, rp.Log)
	case "mail":
		app.sendMail(w, rp.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, l LoggerPayload) {
	//use Marshal() in production....
	jsonData, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		log.Println("Error converting to JSON", err)
		return
	}

	req, err := http.NewRequest("POST", loggerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status NOT OK: %d", resp.StatusCode)
		app.errorJson(w, errors.New("error from Logger Service - status not OK"), resp.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//json to be sent to Auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	
	//call auth microservice
	request, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	log.Printf("The HTTP status for request is: %s", response.Status)

	//check to get the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error when calling AUTH service"))
		return
	}

	//decode the json from the auth service
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
	}

	payload := jsonResponse	{
		Error: false,
		Message: "Authenticated Succesfully Green!",
		Data: jsonFromService.Data,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	jsonData, _ := json.Marshal(m) // should use marshal in production instead of MarshalIndent
	
	// post to Mail-service
	req, err := http.NewRequest("POST", mailerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling mail service"))
		return
	}

	//json response from broker
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + m.To
	
	app.writeJson(w, http.StatusAccepted, payload)
}