package consumer

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Start(ctx context.Context, c *amqp.Channel, cs service.CourseService) error {
	var (
		courseConsumer = NewCourseConsumer(c, cs)
	)

	if err := courseConsumer.Start(ctx); err != nil {
		return err
	}

	return nil
}
