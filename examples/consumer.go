package main

import (
	"flag"
	"log"
	"rmq"
)

func main() {
	url := flag.String("url", "", "url connection to rabbitmq api")
	queue := flag.String("q", "", "name of the queue")

	flag.Parse()

	r, err := rmq.NewRabbitMQConnection(*url)
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	q, err := r.QueueDeclare(rmq.QueueConfig{Name: *queue, Durable: true})
	if err != nil {
		log.Panic(err)
	}
	msgs, err := r.Consume(rmq.ConsumeConfig{Queue: q.Name, Consumer: "", AutoAck: true})
	if err != nil {
		log.Panic(err)
	}

	finish := make(chan struct{})

	go func() {
		for msg := range msgs {
			log.Printf("Received: %s\n", string(msg.Body))
		}
	}()
	log.Println("Waiting for messages... . Press C-X to exit.")
	<-finish
}
