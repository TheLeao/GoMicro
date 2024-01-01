package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// import (
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

type Consumer struct{
	conn *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}

	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Println("Error on declareRandomQueue: ", err)
		return err
	}

	for _, t := range topics {
		err = ch.QueueBind(
			q.Name,
			t,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			log.Println("Error on binding queue: ", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println("Error on consuming messages: ", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(p Payload) {
	// situations based on our own scenario
	// switch p.Name {
	// case "log", "event":
	// 	// log the 
	// 	err := logEvent(p) 
	// 	if err != nil {

	// 	}
	// case "auth":
	// 	// authenticate

	// default:

	// }
}

func logEvent(p Payload) error {
	//use Marshal() in production....
	jsonData, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		log.Println("Error converting to JSON", err)
		return err
	}

	loggerURL := "http://logger-service/log"

	req, err := http.NewRequest("POST", loggerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {		
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error on Log event: status not Ok: ", err)
		return err
	}

	return nil
}