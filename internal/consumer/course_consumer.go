package consumer

import (
	"context"
	"log"

	"github.com/loanem-backend/course-service/pkg/messaging"
	"github.com/loanem-backend/inventory-service/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CourseConsumer struct {
	ch   *amqp.Channel
	serv service.CourseService
}

func NewCourseConsumer(c *amqp.Channel, cs service.CourseService) *CourseConsumer {
	return &CourseConsumer{
		ch:   c,
		serv: cs,
	}
}

func (c *CourseConsumer) Start(ctx context.Context) error {
	queueName := "inventory.course_events.queue"

	q, err := c.ch.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = c.ch.ExchangeDeclare(
		string(messaging.ExchangeCourseEvents),
		"fanout", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(
		q.Name, "",
		string(messaging.ExchangeCourseEvents),
		false, nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			var event messaging.CourseEvent
			if err := event.UnmarshalPayload(&m); err != nil {
				log.Printf("Failed to unmarshal delivery body: %v\n", err)
				m.Nack(false, false)
				continue
			}

			var e error

			switch messaging.EventName(event.EventName) {
			case messaging.EventNameCourseCreated:
				e = c.serv.Add(ctx, event)
			case messaging.EventNameCourseDeleted:
				e = c.serv.Remove(ctx, event)
			default:
				log.Printf("Unrecognized event name: %s\n", event.EventName)
				m.Ack(false)
				continue
			}

			if e != nil {
				log.Printf("Error processing event: %v\n", e)
				m.Nack(false, false)
			} else {
				m.Ack(false)
			}
		}
	}()

	return nil
}
