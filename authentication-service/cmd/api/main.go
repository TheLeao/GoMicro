package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

const webPort = "80" // the same port works inside docker

type Config struct {
	DB *sql.DB
	Models data.Models
}

var (
	dbServer = os.Getenv("DbServer")
	dbPassword = os.Getenv("Password")
	dbUser = os.Getenv("DbUser")
	dbPort = os.Getenv("DbPort")
	database = os.Getenv("Database")
 )

var connString string

func main() {
	log.Println("Starting Auth Service")

	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", dbServer, dbUser, dbPassword, dbPort, database)
	
	conn := connectDB()
	if conn == nil {
		log.Panic("Can't connect to DB!")
	}

	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(srv string) (*sql.DB, error) {
	db, err := sql.Open("mssql", connString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	var counts int

	for {
		conn, err := openDB("")

		if err != nil {
			log.Panicf("Database not ready yet... %s", err)
			counts++
		} else {
			log.Println("Connected to Database!")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for retry...")
		time.Sleep(2 * time.Second)
		continue
	}
}