package config

import (
	"errors"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMqConfig struct {
	URL             string
	Queues          map[string]string
	RabbitMqChannel *amqp091.Channel
}

func NewRabbitMqConfig() (*RabbitMqConfig, error) {
	url := os.Getenv("RABBITMQ_URL")
	if len(url) == 0 {
		return nil, errors.New("environment variable RABBITMQ_URL not initialized")
	}

	queueEventNotifications := os.Getenv("RABBITMQ_QUEUE_EVENT_NOTIFICATIONS")
	if len(queueEventNotifications) == 0 {
		return nil, errors.New("environment variable RABBITMQ_QUEUE_EVENT_NOTIFICATIONS not initialized")
	}

	queueBroadcastNotifications := os.Getenv("RABBITMQ_QUEUE_BROADCAST_NOTIFICATIONS")
	if len(queueBroadcastNotifications) == 0 {
		return nil, errors.New("environment variable RABBITMQ_QUEUE_EVENT_NOTIFICATIONS not initialized")
	}

	rqc := RabbitMqConfig{
		URL: url,
		Queues: map[string]string{
			"event_notifications":     queueEventNotifications,
			"broadcast_notifications": queueBroadcastNotifications,
		},
	}

	if err := rqc.connectRabbitMq(); err != nil {
		return nil, err
	}

	return &rqc, nil
}

func (a *RabbitMqConfig) connectRabbitMq() error {
	conn, err := amqp091.Dial(a.URL)
	if err != nil {
		log.Fatalf("Failed to connect to rabbit mq: %s", err)
	}
	log.Print("Connect to rabbit mq")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel rabbit mq: %s", err)
	}
	log.Print("Open a channel rabbit mq")

	for _, queueName := range a.Queues {
		_, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Fatalf("failed to declare queue: %s", err)
		}
		log.Printf("Declate queue: %s", queueName)
	}

	a.RabbitMqChannel = channel
	log.Print("Set rabbit mq channel")

	return nil
}
