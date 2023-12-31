package main

import (
	"log"
	"net/http"
)

//const webPort string = "80"

type Config struct{}

func main() {
	app := Config{}	

	// log.Printf("Starting Broker Service on port %s\n", webPort)
	log.Printf("Starting Broker Service")

	// define http server
	srv := &http.Server {
		Addr: ":80",
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}