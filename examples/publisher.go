package main

import (
	"encoding/json"
	"flag"
	"log"
	"rmq"

	"github.com/streadway/amqp"
)

type BodyS struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	url := flag.String("url", "", "url connection to rabbitmq api")
	exchange := flag.String("exchange", "", "exchange's name")
	key := flag.String("key", "", "routing key. If exchange empty key is the queue's name")

	flag.Parse()

	r, err := rmq.NewRabbitMQConnection(*url)
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	for i := 0; i < 1000; i++ {
		body := BodyS{
			Id:   i,
			Name: "coca cola zero",
		}
		raw, err := json.Marshal(body)
		if err != nil {
			log.Panic(err)
		}
		if err := r.Publish(rmq.PublishConfig{Exchange: *exchange, Key: *key, Msg: amqp.Publishing{ContentType: "plain/text", Body: raw}}); err != nil {
			log.Println(err)
		}
	}
}
