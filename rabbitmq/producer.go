package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/model"
	"github.com/streadway/amqp"
)

// Producer produce data into RabbitMQ
type Producer struct {
	AMPQChannel *ChannelWrapper
	Conn        *ConnectionWrapper
	queueName   string
}

// NewProducer create new instance of RabbitMQ producer
func NewProducer(cfg config.Rabbitmq) *Producer {
	c := CreateConnection(cfg)
	p := &Producer{
		AMPQChannel: c.Chann,
		Conn:        c.Conn,
		queueName:   c.QueueName,
	}

	return p
}

func (p *Producer) Queue(d model.Data) error {
	msg, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal json %w", err)
	}

	if err := p.AMPQChannel.Channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		}); err != nil {
		return fmt.Errorf("failed to publish on queue %w", err)
	}

	return nil
}
