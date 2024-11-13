package external

import (
	"AuthService/internal/config"
	"AuthService/internal/convertor"
	"AuthService/internal/entity"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const (
	eventQueueName = ""
)

type NotificationExternal struct {
	rabbitMqConfig *config.RabbitMqConfig
}

func NewNotificationExternal(rq *config.RabbitMqConfig) *NotificationExternal {
	return &NotificationExternal{
		rabbitMqConfig: rq,
	}
}

func (ne *NotificationExternal) SendEmailEventNotification(en *entity.EmailEventNotificationEntity) error {
	req := convertor.EmailEventNotificationEntityToProto(en)
	body, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	return ne.rabbitMqConfig.RabbitMqChannel.Publish(
		"",
		eventQueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/protobuf",
			Body:        body,
		},
	)
}

func (ne *NotificationExternal) SendEmailBroadcastNotification() {}
