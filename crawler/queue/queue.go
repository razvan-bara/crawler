package queue

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

var (
	ErrorQueueMessageDuplicate = errors.New("queue message duplicate")
)

type MessageQueue struct {
	name  string
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

func NewMessageQueue(conn *amqp.Connection, queueName string) (*MessageQueue, error) {
	queue, channel, err := getQueue(conn, queueName)
	if err != nil {
		return nil, fmt.Errorf("error getting queue %s: %w", queueName, err)
	}

	return &MessageQueue{
		conn:  conn,
		ch:    channel,
		queue: queue,
		name:  queueName,
	}, nil
}

func (queue *MessageQueue) Publish(msg amqp.Publishing) error {
	if !isUniqueMessage(msg.Body) {
		return ErrorQueueMessageDuplicate
	}
	err := queue.ch.Publish(
		"",
		queue.name,
		false,
		false,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func (queue *MessageQueue) GetConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := queue.ch.Consume(
		queue.name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return msgs, nil
}

func ConnectToRabbitMQ(user, password, host, port, queueVirtualHost string) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, queueVirtualHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getQueue(conn *amqp.Connection, queueName string) (*amqp.Queue, *amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &q, ch, nil
}
