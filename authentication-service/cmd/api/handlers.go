package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type requestPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	log.Println("AUTHENTICATING....")

	rp := &requestPayload{}

	err := app.readJson(w, r, rp)
	if err != nil {
		log.Println("Read JSON failed...", err)
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	//validate user against the database
	user, err := app.Models.User.GetByEmail(rp.Email)
	if err != nil {
		log.Println("Get Email Failed....", err)
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(rp.Password)
	if !valid || err != nil {
		log.Println("Get Password Matches failed....", err)
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// logging authentication in Logger Service
	err = app.logRequest("Authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
	}

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in as %s", user.Email),
		Data: user,
	}

	log.Println("Authentication process completed.")
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log" //as in docker-compose

	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {	
		return fmt.Errorf("error creating log request: %s", err)
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error on log request: %s", err)
	}

	return nil
}