package amqpman

import (
	"context"
	"github.com/betam/glb/lib2"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func Publish(publisherCh *amqp.Channel, context context.Context, exchange string, routingKey string, payload []byte) {

	err := publisherCh.PublishWithContext(
		context,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		},
	)
	lib2.LogOnError(err, "Failed to publish a message", "warn")
	log.Tracef("[R>] <%v>", routingKey)

}
