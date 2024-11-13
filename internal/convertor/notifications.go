package convertor

import (
	"AuthService/internal/entity"

	notifiocationProto "github.com/sergeyiksanov/notification-service/pkg/api/v1"
)

func EmailEventNotificationEntityToProto(en *entity.EmailEventNotificationEntity) *notifiocationProto.EventNotificationRequest {
	return &notifiocationProto.EventNotificationRequest{
		Email: en.Email,
		Name:  en.Name,
		Title: en.Title,
		Body:  en.Body,
	}
}
