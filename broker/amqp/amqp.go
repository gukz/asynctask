package amqp

import (
	"crypto/tls"
	"errors"
	"github.com/gukz/asynctask"
	"github.com/streadway/amqp"
	"time"
)

var _ asynctask.Broker = (*amqpBroker)(nil)

type amqpBroker struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	tlsConfig   *tls.Config
	queue       amqp.Queue
	confirmChan chan amqp.Confirmation

	url             string
	exchange        string
	queueName       string
	queueBindingKey string
}

func NewBroker(url string, exchange string, queueName string, queueBindingKey string) (asynctask.Broker, error) {
	a := &amqpBroker{
		url:             url,
		exchange:        exchange,
		queueName:       queueName,
		queueBindingKey: queueBindingKey,
	}
	var err error

	tlsConfig := &tls.Config{}
	a.conn, err = amqp.DialTLS(url, tlsConfig)
	if err != nil {
		return nil, err
	}
	if a.channel, err = a.conn.Channel(); err != nil {
		return nil, err
	}
	if err := a.channel.ExchangeDeclare(
		a.exchange, // name of the exchange
		"topic",    // type
		true,       // durable
		false,      // delete when complete
		false,      // internal
		false,      // noWait
		amqp.Table{
			"x-message-ttl": 1 * time.Hour,
			"x-expires":     1 * time.Hour,
		}, // arguments
	); err != nil {
		return nil, err
	}
	if a.queue, err = a.channel.QueueDeclare(
		a.queueName, // name
		false,       // durable
		true,        // delete when unused
		false,       // exclusive
		false,       // no-wait
		amqp.Table{ // arguments
			"x-message-ttl": 1 * time.Hour,
			"x-expires":     1 * time.Hour,
		},
	); err != nil {
		return nil, err
	}
	if err = a.channel.QueueBind(
		a.queue.Name,      // name of the queue
		a.queueBindingKey, // binding key
		a.exchange,        // source exchange
		false,             // noWait
		amqp.Table{},      // arguments
	); err != nil {
		return nil, err
	}
	if err := a.channel.Confirm(false); err != nil {
		return nil, err
	}
	a.confirmChan = a.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	return a, nil
}

func (a *amqpBroker) CheckHealth() bool {
	return true
}

func (a *amqpBroker) AckMessage(queue string, messageId uint64) error {
	a.channel.Ack(messageId, false)
	return nil
}

func (a *amqpBroker) PopMessage(queue string) (*asynctask.BrokerMessage, error) {
	d, ok, err := a.channel.Get(queue, false)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &asynctask.BrokerMessage{Body: d.Body, Id: d.DeliveryTag}, nil
}

func (a *amqpBroker) PushMessage(queue string, data []byte) error {
	if err := a.channel.Publish(
		a.exchange,  // exchange
		a.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // Persistent // Transient
		},
	); err != nil {
		return err
	}
	confirmed := <-a.confirmChan
	if confirmed.Ack {
		return nil
	}
	return errors.New("Error happen when publish message.")
}

func (a *amqpBroker) Close() error {
	return a.channel.Close()
}
