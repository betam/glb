package amqpman

import (
	"fmt"
	"github.com/betam/glb/lib2"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateExchange(ch *amqp.Channel, name string, panic bool, ttl int32) error {

	args := make(amqp.Table)
	if ttl > 0 {
		args["x-message-ttl"] = ttl
	}

	err := lib2.AmqpDeclareExchange(
		ch, name,
		"topic",
		true,
		false,
		false,
		false,
		args,
	)

	if panic {
		lib2.LogOnError(err, fmt.Sprintf("Failed while exchange create (%v)", name), "panic")
	}

	return err
}

func BindToExchange(ch *amqp.Channel, queueName string, routingKey string, exchangeName string, panic bool) error {
	err := lib2.AmqpQueueBind(
		ch,
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if panic {
		lib2.LogOnError(err, fmt.Sprintf("Failed while binding pipe to exchangee create (%v, %v, %v)", queueName, routingKey, exchangeName), "panic")
	}

	return err
}
