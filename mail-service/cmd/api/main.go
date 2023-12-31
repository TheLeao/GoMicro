package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct{
	Mailer Mail
}

const webPort = "80"

func main() {

	app := Config{
		Mailer: newMailer(),
	}

	log.Println("Starting Mail Service at port: ", webPort)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func newMailer() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return Mail {
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("FROM_NAME"),
		FromAddres: os.Getenv("FROM_ADDRESS"),
	}
}