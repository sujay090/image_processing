package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect establishes a connection to RabbitMQ
func Connect(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close gracefully closes the connection and channel
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// DeclareQueue creates a queue if it doesn't exist
func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// Publish sends a message to the specified queue
func (r *RabbitMQ) Publish(ctx context.Context, queueName string, body []byte) error {
	// Add a timeout for the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	return err
}

// Consume returns a channel of deliveries from the specified queue
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName, // queue
		"",        // consumer tag
		false,     // auto-ack (we will manually ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}
