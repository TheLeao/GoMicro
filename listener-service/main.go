package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" // alias used to make reference cleaner
)

func main() {
	// connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// listen to messages
	log.Println("Listening for and consuming messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println("Error when listening", err)
	}
}

func connect() (*amqp.Connection, error) {
	// RabbitMQ may take a while to start
	var counts int64 = 0
	var conn *amqp.Connection
	backOff := 1 * time.Second

	// waiting for rabbit to be ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost") 
		if err != nil {
			fmt.Println("RabbitMQ not ready yet")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			conn = c
			break
		}

		if counts > 5 {
			fmt.Print("Timeout error: ", err)
			return nil, err
		}

		// exponentially increasing the time between tries before timeout
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off ", counts)
		time.Sleep(backOff)
	}

	return conn, nil
}