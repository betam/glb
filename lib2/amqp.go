package lib2

import (
	xc "github.com/shomali11/util/xconditions"
	log "github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

func debugOrPanic(panic bool) string {
	return xc.IfThenElse(panic, "panic", "debug").(string)
}

func AmqpConnect(url string, panic bool) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)

	LogOnError(err, "Failed to connect to RabbitMQ", debugOrPanic(panic))
	log.Debugf("Amqp connected: %v", url)
	return conn, err
}

func AmqpOpenChanel(conn *amqp.Connection, panic bool) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	LogOnError(err, "Failed to open a chanel", debugOrPanic(panic))
	return ch, err
}

func AmqpConnectAndChanel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, _ := AmqpConnect(url, true)
	ch, _ := AmqpOpenChanel(conn, true)
	return conn, ch
}

func AmqpConsumer(ch *amqp.Channel, qName string, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(qName, consumer, autoAck, exclusive, noLocal, noWait, args)
	LogOnError(err, "Failed to register a consumer", "debug")

	log.Debugf("Consumer started, queue: %v", qName)
	return msgs, err
}

func AmqpConsumerLite(ch *amqp.Channel, qName string, autoAck bool, exclusive bool) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(qName, "", autoAck, exclusive, false, false, nil)
	LogOnError(err, "Failed to register a consumer", "panic")

	log.Debugf("Consumer started, queue: %v", qName)
	return msgs, err
}

func AmqpDeclareQueue(ch *amqp.Channel, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	LogOnError(err, "Failed to declare a queue", "debug")
	log.Debugf("Declared  queue: %v", name)
	return q, err
}

func AmqpDeclareExchange(ch *amqp.Channel, name string, etype string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error {
	err := ch.ExchangeDeclare(name, etype, durable, autoDelete, internal, noWait, args)
	LogOnError(err, "Failed to declare a exchange", "debug")
	log.Debugf("Declared  exchange: %v", name)
	return err
}

func AmqpQueueBind(ch *amqp.Channel, name string, key string, exchange string, noWait bool, args amqp.Table) error {
	err := ch.QueueBind(name, key, exchange, noWait, args)
	LogOnError(err, "Failed to bind a queue to exchange", "debug")
	log.Debugf("Queue  binding: %v-queue to %v-exchange via %v-key ", name, exchange, key)
	return err
}
