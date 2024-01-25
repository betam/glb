package amqpman

import (
	"fmt"
	"github.com/betam/glb/lib2"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateQueue(ch *amqp.Channel, name string, panic bool, ttl int32) (amqp.Queue, error) {

	args := make(amqp.Table)
	if ttl > 0 {
		args["x-message-ttl"] = ttl
	}

	q, err := lib2.AmqpDeclareQueue(
		ch, name,
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		args,
	)

	if panic {
		lib2.LogOnError(err, fmt.Sprintf("Failed while queue create (%v)", name), "panic")
	}

	return q, err
}
