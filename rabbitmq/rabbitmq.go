package rabbitmq

import (
	"fmt"
	"time"

	"github.com/I1820/I1820/config"
	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

const reconnectLimit = 300
const reconnectPeriod = 1 * time.Second

// MessageTTL provides Time to Live for each message in Queue (Milliseconds)
const MessageTTL = 60 * 1000

// ConnectionWrapper wraps AMQP Connection
type ConnectionWrapper struct {
	Connection *amqp.Connection
}

// ChannelWrapper wraps AMQP Channel
type ChannelWrapper struct {
	Channel *amqp.Channel
}

// Connection to the RabbitMQ
type Connection struct {
	Conn      *ConnectionWrapper
	Chann     *ChannelWrapper
	QueueName string
}

// CreateConnection creates a RabbitMQ connection
func CreateConnection(cfg config.Rabbitmq) *Connection {
	name := fmt.Sprintf("%s-%s", config.Namespace, "sms")
	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
	)

	conn, channel := connect(addr, name)

	return &Connection{
		Conn:      conn,
		Chann:     channel,
		QueueName: name,
	}
}

func connect(addr string, queueName string) (*ConnectionWrapper, *ChannelWrapper) {
	conn := createConnection(addr)

	channel := createChannel(conn, queueName)

	return conn, channel
}

func createConnection(addr string) *ConnectionWrapper {
	conn := connectionWithRetry(addr)

	connectionWrapper := &ConnectionWrapper{Connection: conn}

	// watch for connection failures
	go func() {
		for {
			reason, ok := <-connectionWrapper.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				logrus.Info("rabbit: connection closed without error.")
				break
			}

			logrus.Errorf("rabbit: connection closed. error=%s", reason.Error())

			time.Sleep(reconnectPeriod)

			conn := connectionWithRetry(addr)

			connectionWrapper.Connection = conn

			logrus.Info("rabbit: connection reconnect success")
		}
	}()

	return connectionWrapper
}

func connectionWithRetry(addr string) *amqp.Connection {
	for i := 0; i < reconnectLimit; i++ {
		conn, err := amqp.Dial(addr)
		if err != nil {
			logrus.Errorf("Failed to connect to Rabbitmq at %s: (%s). Retrying... %d ", addr, err.Error(), i)
			time.Sleep(reconnectPeriod)
		} else {
			return conn
		}
	}
	panic("rabbit: connection retry limit reached.")
}

func createChannel(cw *ConnectionWrapper, queueName string) *ChannelWrapper {
	channel := channelWithRetry(cw, queueName)

	channelWrapper := &ChannelWrapper{Channel: channel}

	go func() {
		for {
			reason, ok := <-channelWrapper.Channel.NotifyClose(make(chan *amqp.Error))
			if !ok {
				logrus.Info("rabbit: channel closed without error.")
				break
			}

			logrus.Errorf("rabbit: channel closed. error=%s", reason.Error())

			time.Sleep(reconnectPeriod)

			channel := channelWithRetry(cw, queueName)

			channelWrapper.Channel = channel

			logrus.Info("rabbit: channel reconnect success")
		}
	}()

	return channelWrapper
}

func channelWithRetry(cw *ConnectionWrapper, queueName string) *amqp.Channel {
	for i := 0; i < reconnectLimit; i++ {
		ch, err := cw.Connection.Channel()
		if err != nil {
			logrus.Errorf("Failed to create channel (%s). Retrying... %d ", err.Error(), i)

			time.Sleep(reconnectPeriod)

			continue
		}

		const (
			prefetchCount = 25
			prefetchSize  = 0
		)

		if err := ch.Qos(prefetchCount, prefetchSize, false); err != nil {
			logrus.Errorf("Failed to configure Qos for channel (%s). Retrying... %d ", err.Error(), i)

			time.Sleep(reconnectPeriod)

			continue
		}

		if _, err := ch.QueueDeclare(
			queueName, true, false, false, false, amqp.Table{
				"x-message-ttl": MessageTTL,
			},
		); err != nil {
			logrus.Errorf("Failed to declare queue for channel (%s). Retrying... %d ", err.Error(), i)

			time.Sleep(reconnectPeriod)

			continue
		}

		return ch
	}
	panic("rabbit: channel creation limit reached.")
}
