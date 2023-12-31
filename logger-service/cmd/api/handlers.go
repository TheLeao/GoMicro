package main

import (
	"log"
	"logger-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload jsonPayload
	_ = app.readJson(w, r, &requestPayload)

	// data to insert
	event := data.LogEntry {
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Println("Error logging Log", err)
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse {
		Error: false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusOK, resp)
}