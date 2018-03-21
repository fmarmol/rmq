package rmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type QueueBinding struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}

type PublishConfig struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Msg       amqp.Publishing
}

type ConsumeConfig struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type ExchangeConfig struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type RabbitMQConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (r *RabbitMQConnection) Close() {
	r.ch.Close()
	r.conn.Close()
}

func (r *RabbitMQConnection) QueueDeclare(cfg QueueConfig) (amqp.Queue, error) {
	q, err := r.ch.QueueDeclare(cfg.Name, cfg.Durable, cfg.AutoDelete, cfg.Exclusive, cfg.NoWait, cfg.Args)
	return q, err
}

func (r *RabbitMQConnection) QueueBind(cfg QueueBinding) error {
	return r.ch.QueueBind(cfg.Name, cfg.Key, cfg.Exchange, cfg.NoWait, cfg.Args)
}

func (r *RabbitMQConnection) Publish(cfg PublishConfig) error {
	return r.ch.Publish(cfg.Exchange, cfg.Key, cfg.Mandatory, cfg.Immediate, cfg.Msg)
}

func (r *RabbitMQConnection) Consume(cfg ConsumeConfig) (<-chan amqp.Delivery, error) {
	msg, err := r.ch.Consume(cfg.Queue, cfg.Consumer, cfg.AutoAck, cfg.Exclusive, cfg.NoLocal, cfg.NoWait, cfg.Args)
	return msg, err
}

func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create channel: %v", err)
	}
	return &RabbitMQConnection{conn: conn, ch: ch}, nil
}
