package mappers

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
)

func ToDB(message *models.Message) *db.Chat {

	return &db.Chat{
		Sender:   message.Sender,
		Body:     message.Body,
		Reciever: message.Reciever,
		Sent:     uint64(message.Sent),
	}
}

func ToNetwork(chat *db.Chat) *models.Message {
	return &models.Message{
		Sender:   chat.Sender,
		Body:     chat.Body,
		Reciever: chat.Reciever,
		Sent:     int64(chat.Sent),
	}
}
